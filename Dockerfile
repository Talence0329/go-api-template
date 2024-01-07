FROM centos:7
MAINTAINER TEST 

ARG filename
ARG version

ENV File=${filename}_${version}
ADD ${File} ${File}
ADD run.sh run.sh 
ADD resource resource
RUN chmod +x /${File}
RUN chmod +x /run.sh 
ENTRYPOINT ["/run.sh"]
