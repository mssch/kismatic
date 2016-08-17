FROM centos:7

RUN yum -y update && yum clean all

RUN yum -y install python-setuptools gcc python-devel libffi-devel openssl-devel

RUN easy_install pip 

RUN pip install --upgrade setuptools --user python
