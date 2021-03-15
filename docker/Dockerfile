FROM ubuntu
MAINTAINER youht

WORKDIR /opt

RUN apt-get update -y
RUN apt-get install -y mongodb
RUN apt-get install -y mysql-client 

#add alpine-sdk when alpine
 
RUN apt-get install -y  \
            bash \
            net-tools \
            iputils-ping \
            wget \
            curl \
            tree \
            vim \
            tmux \
            netcat-traditional          
#add openssl-dev when alpine
RUN apt-get install -y gcc g++  

RUN apt-get install -y python3 \
            python3-dev \
            python3-pip \
            python-dev \
            python-pip 
#    && pip3 install --upgrade pip \
#    && pip install --upgrade pip
    
#RUN pip3 install torch \
#             torchvision \
#             torchtext

#add zeromq zeromq-dev
RUN apt-get install -y libczmq-dev \
    && pip3 install jupyter



RUN pip3 install numpy

#add libxslt when alpine
RUN pip3 install pandas \
    && apt-get install -y libxml2 \
               libxml2-dev \
               libxslt-dev  \
    && pip3 install lxml \
                    html5lib \
                    beautifulsoup4 \
                    xlrd

#add jpeg-dev zlib-dev
RUN apt-get install -y libjpeg-dev \
        zlib1g-dev \
    && pip3 install pillow

#add libpng freetype freetype-dev when alpine 
RUN apt-get install -y libpng-dev \
    && pip3 install matplotlib


RUN python2 -m pip install ipykernel \
    && python2 -m ipykernel install --user

    
RUN pip3 install jupyter_contrib_nbextensions \
         jupyter_nbextensions_configurator \
    && pip3 install jupyterthemes \
    && jupyter contrib nbextension install --user \
    && jupyter nbextensions_configurator enable --user



#python package
RUN pip3 install findspark \
   flask flask-socketio flask-cors \
   requests bash pyyaml \
   pytz socketio-client \
   pycrypto \
   pymongo bson
                 
#node_modules install base /root

RUN pip3 install opencv-python \
                 Progressbar \
                 pymysql \
                 psutil \
                 blaze \
                 pyecharts \
                 tushare \
                 kafka \
                 jieba \
                 scipy


RUN apt-get install -y clinfo \
    && pip3 install pybind11 \
                    mako \
                    pyopencl

RUN pip3 install autopep8


#patch 2019.07.02
RUN pip3 install gensim


COPY /*.sh /
RUN chmod 700 /*.sh

#RUN sed -i 's/"${PYSPARK_PYTHON:-"python"}"/"${PYSPARK_PYTHON:-"ipython3"}"/' /opt/spark-2.4.3/bin/pyspark
ENV PYSPARK_DRIVER_PYTHON=ipython3
ENV PYSPARK_PYTHON=ipython3
ENV PYTHONIOENCODING="utf8"

ENV HOME /root
ENV SERVER_PORT 8888
ENV FLASK_PORT  3000
ENV NODEJS_PORT 4000
ENV JUPYTER_HUB_PORT 8000

EXPOSE $SERVER_PORT \
       $JUPYTER_HUB_PORT \
       $FLASK_PORT \
       $NODEJS_PORT

COPY /custom.js /root/.jupyter/custom/

#install jupyterhub
RUN pip3 install jupyterhub \
    && npm install -g configurable-http-proxy

#patch local user -> system-wide scope
RUN jupyter contrib nbextension install --system
RUN /opt/./almond --install --global
RUN cd /opt/IJava \
    && mkdir -p /usr/local/share/jupyter/kernels/java \
    && ./gradlew installKernel --path /usr/local/share/jupyter/
RUN cd /opt/jupyter-nodejs \
    && mkdir -p /usr/local/share/jupyter/kernels/nodejs/ \
    && node install.js /usr/local/share/jupyter/kernels/nodejs/

# 2019.07.22 install go & igo
ADD go /opt/go
ENV GOROOT=/opt/go
ENV GOPATH=/root/source/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

RUN apt-get install pkg-config
#RUN go get -u github.com/gopherdata/gophernotes
RUN mkdir -p $GOPATH/bin 
COPY /gokernel/bin/gophernotes /opt/go/bin/
RUN mkdir -p ~/.local/share/jupyter/kernels/gophernotes
COPY /gokernel/kernel/* /root/.local/share/jupyter/kernels/gophernotes/

#RUN cp /go/src/github.com/gopherdata/gophernotes/kernel/* ~/.local/share/jupyter/kernels/gophernotes

# 2019.08.01 use mod & get some module
ENV GO111MODULE=auto
ENV GOPROXY=https://goproxy.io

RUN go get -u -v github.com/segmentio/kafka-go
RUN go get -u -v github.com/go-sql-driver/mysql
#RUN go get -u -v github.com/kataras/iris
RUN go get -u -v github.com/astaxie/beego
RUN go get -u -v github.com/beego/bee
RUN go get -u -v gopkg.in/mgo.v2
RUN go install gopkg.in/mgo.v2/bson
RUN go get -u -v github.com/gorilla/mux
RUN go get -u -v github.com/codegangsta/negroni
RUN go get -u -v github.com/unrolled/render

RUN mkdir -p $GOPATH/src/golang.org/x \
    && cd $GOPATH/src/golang.org/x  \
    && git clone https://github.com/golang/net.git  
RUN go get -u -v github.com/levigross/grequests
#RUN go get -u -v go.mongodb.org/mongo-driver

#RUN go get -v github.com/tjfoc/gmsm
ADD tjfoc $GOPATH/src/github.com/tjfoc/
RUN go install github.com/tjfoc/gmsm/sm2 \
    && go install github.com/tjfoc/gmsm/sm3 \
    && go install github.com/tjfoc/gmsm/sm4
#RUN go get -u -v github.com/tidwall/gjson
ADD tidwall $GOPATH/src/github.com/tidwall/
RUN go install github.com/tidwall/gjson

RUN echo "export JAVA_HOME=/opt/jdk-10.0.2 \
          export SPARK_HOME=/opt/spark-2.4.3 \
          export PATH=$PATH:$JAVA_HOME/bin:$SPARK_HOME/bin:$SPARK_HOME/sbin \
          export PYSPARK_DRIVER_PYTHON=ipython3 \
          export PYSPARK_PYTHON=ipython3 \
          export PYTHONIOENCODING=utf8" > /etc/profile

RUN echo "nohup jupyter notebook --ip=0.0.0.0 --port=$SERVER_PORT --notebook-dir=$HOME --allow-root >jupyter.log 2>&1 &" > /root/start-jupyter.sh && \
    echo "nohup jupyterhub &" >> /root/start-jupyter.sh && \
    echo "tail -f jupyter.log" >> /root/start-jupyter.sh && \
    chmod 700 /root/start-jupyter.sh
CMD /root/start-jupyter.sh
