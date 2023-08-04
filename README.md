# 2PCTM

This is a demo version of 2 transaction manager with 3 demo services

basic idea:
There are 2 phase for the operation
1. tm_server will enter prepare phase, where will call each service prepare
2. tm_server will enter commit phase, where will call each service commit 