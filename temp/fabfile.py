# -*- coding: utf-8 -*-

from fabric.api import *
from fabric.decorators import parallel
from fabric.contrib.files import *
import time
import os.path

REPO_URL = "ssh://git@xxxx/sqssendmessage.git"
DOCKER_REPO_URI_BASE = "aws-ecr/sqssendmessage"

env.use_ssh_config = True
env.forward_agent = True

def development(*hosts):
    """ Development environment """
    env.environment = "development"
    env.exec_user = "sally"
    env.hosts = hosts
    env.target_branch = "develop"

def staging(*hosts):
    """ Staging environment """
    env.environment = "staging"
    env.exec_user = "sally"
    env.hosts = hosts
    env.target_branch = "staging"

def production(*hosts):
    """ Production environment """
    env.environment = "production"
    env.exec_user = "sally"
    env.hosts = hosts
    env.target_branch = "master"

def _rewrite_config():
    # Aligning target branches of GitLab repository
    local("sed -i -e 's/\(feature.*\|develop\|staging\|master\)/%s/g' ./glide.yaml" % env.target_branch, shell='/bin/bash')


@runs_once
def build():
    composefile = "docker-compose.%s.yml" % env.environment
    _updateLocalCredential()
    _updateLibrary()    
    _buildLatestDockerImage(composefile)


@parallel
def deploy():
    """
    Stop and deploy sqssendmessage, then start it.
    """
    composefile = "docker-compose.%s.yml" % env.environment
    _updateRemoteCredential()
    remotepath = "/home/xxxx/sqssendmessage/"
    run('mkdir -p %s', remotepath)

    with lcd('config'):
        sudo('mkdir -p /etc/sqssendmessage')
        with cd('/etc/sqssendmessage'):
            put('%s.toml' % env.environment, 'sqssendmessage.toml', use_sudo=True, mode=0666)

    with cd(remotepath):
        _updateRemoteDockerImage(composefile, remotepath)
        _restartContainer(composefile)

    local('git checkout ./glide.yaml', shell='/bin/bash')


@runs_once
# glide環境もコンテナで揃えたいが、コンテナ内からプライベートリポジトリにアクセスするのが困難なため
def _updateLibrary():
    _rewrite_config()
    local('go get github.com/Masterminds/glide', shell='/bin/bash')
    local('go install github.com/Masterminds/glide', shell='/bin/bash')
    local('glide up', shell='/bin/bash')

    # inspections
    local('find ./ -name *.go | grep -v vendor | xargs -L1 go vet', shell='/bin/bash')
    local('go test $(glide novendor)', shell='/bin/bash')

@runs_once
def _buildLatestDockerImage(composefile):
    """
    Build new image and push to aws ecr
    """
    local("docker-compose -f %s build" % composefile)
    local("docker-compose -f %s push" % composefile)

@parallel
def _updateRemoteDockerImage(composefile, remotepath):
    put(composefile, remotepath)
    run("docker-compose -f %s pull" % composefile)

@runs_once
def _updateLocalCredential():
    local("eval $(aws ecr get-login --region ap-northeast-1)", shell="/bin/bash")

@parallel
def _updateRemoteCredential():
    run("eval $(aws ecr get-login --region ap-northeast-1)", shell="/bin/bash")


def _restartContainer(composefile):
    run("docker-compose -f %s down" % composefile)
    run("docker-compose -f %s up -d" % composefile)

