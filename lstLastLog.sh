#!/bin/bash
# cd ~/MyProject/Showroom/SRCGI
flast=`ls -rt *.txt | tail -n 1`
less +F $flast
