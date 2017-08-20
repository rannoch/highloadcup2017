# Наследуемся от CentOS 7
FROM centos

USER root
# Выбираем рабочую папку
WORKDIR /root

# Устанавливаем wget и скачиваем Go
RUN yum install -y wget && \
    wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz

RUN rpm -U http://opensource.wandisco.com/centos/7/git/x86_64/wandisco-git-release-7-2.noarch.rpm \
    && yum install -y git

# Устанавливаем Go, создаем workspace и папку проекта
RUN tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz && \
    mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg && \
    mkdir go/src/highloadcup2017
RUN mkdir data

RUN yum install zip & yum install unzip -y

RUN wget http://repo.mysql.com/mysql-community-release-el7-5.noarch.rpm
RUN rpm -ivh mysql-community-release-el7-5.noarch.rpm
RUN yum update -y
RUN yum install mysql-server -y
#
RUN mysql_install_db
#COPY storage/my.cnf /etc/

# Задаем переменные окружения для работы Go
ENV PATH=${PATH}:/usr/local/go/bin GOROOT=/usr/local/go GOPATH=/root/go

# Копируем наш исходный main.go внутрь контейнера, в папку go/src/dumb
ADD . go/src/highloadcup2017

# Компилируем и устанавливаем наш сервер

RUN go get highloadcup2017 && go build highloadcup2017 && go install highloadcup2017

# Открываем 80-й порт наружу
EXPOSE 80

# Запускаем наш сервер
COPY run.sh go/src/highloadcup2017
CMD sh go/src/highloadcup2017/run.sh