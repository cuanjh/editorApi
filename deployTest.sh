#创建目录

echo "开始编译"
obj=editorAPILinux

GOOS=linux GOARCH=amd64 go build -mod=vendor -o $obj

tarFile=${obj}.tar.gz

echo "压缩文件"
tar -cvzf ${tarFile} supervisor-dev.conf ${obj}  static
echo "同步文件"
scp  ${tarFile} root@dev.api.talkmate.com:/tmp/

ssh root@dev.api.talkmate.com >/dev/null 2>&1 << eeooff

exeObj=editorAPILinux

proDir=/opt/data/goPro/\${exeObj}
if [ ! -d \${proDir} ];then
	mkdir -p \${proDir}
fi


mv /tmp/\${exeObj}.tar.gz \${proDir}

cd \${proDir}
tar -xvf \${exeObj}.tar.gz

rm -f /opt/data/goPro/\${exeObj}/static/config/config.json

ln -s /opt/data/goPro/\${exeObj}/static/config/config-dev.json /opt/data/goPro/\${exeObj}/static/config/config.json

if [ ! -L /usr/local/bin/\${exeObj} ];then
        ln -s \${proDir}/\${exeObj} /usr/local/bin
fi

if [ ! -L /usr/local/etc/supervisor.d/\${exeObj}.ini ];then
        ln -s \${proDir}/supervisor-dev.conf /usr/local/etc/supervisor.d/\${exeObj}.ini
fi

supervisorctl update
supervisorctl restart \${exeObj}


exit
eeooff

rm -f ${tarFile}

echo done!
