
obj=editorAPILinux

host=182.92.117.235

GOOS=linux GOARCH=amd64 go build  -o $obj

rsync -avz static root@$host:/opt/data/goPro/editorApi

scp  $obj root@$host:/tmp/

ssh root@$host > /dev/null 2>&1 << eeooff

rm -f /opt/data/goPro/editorApi/static/config/config.json

ln -s /opt/data/goPro/editorApi/static/config/config-prod.json /opt/data/goPro/editorApi/static/config/config.json

ps aux|grep ./editorAPILinux|grep -v grep|grep -v bash|awk '{print \$2}'|xargs kill -9

mv /tmp/editorAPILinux /opt/data/goPro/editorApi/

cd /opt/data/goPro/editorApi/

export export GIN_MODE=release

nohup ./editorAPILinux >>nohup.out 2>&1 &

ps aux|grep ./editorAPI|grep -v grep

exit
eeooff
echo done!
