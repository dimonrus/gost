printenv | sed 's/^\(.*\)$/export \1/g' > variables.sh
chmod 0777 variables.sh
mkfifo -m 0666 /var/log/cron.log
crond -f && tail -f /var/log/cron.log