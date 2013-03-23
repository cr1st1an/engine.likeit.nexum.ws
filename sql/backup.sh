#!/bin/bash
FILE=$USER-$(date +%Y%m%d%H%M%S.%N).sql
mysqldump -h dev-likeit-db.ckmmqpiq5icd.us-east-1.rds.amazonaws.com -u dev_likeit_user dev_likeit_data -p123456 > $FILE
gzip $FILE
