FROM amazon/dynamodb-local:latest
USER root
RUN yum install -y zip
COPY log4j2.xml log4j2.xml
RUN zip -d DynamoDBLocal.jar log4j2.xml && zip -u DynamoDBLocal.jar log4j2.xml