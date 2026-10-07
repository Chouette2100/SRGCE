#!/bin/sh
cd /var/lib/srgce
./SRGCE > /dev/null 2>> error.log
date >> error.log
