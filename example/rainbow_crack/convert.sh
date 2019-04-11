#!/bin/bash
FILE_NAME=$1

cat $FILE_NAME | sort | uniq -f1 > converted_table.txt