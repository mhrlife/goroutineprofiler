# Goroutine Profiler

### Story
We have had a goroutine leak in our system which used to happen completely random without any pattern.We wanted to profiler our system but the problem was we didn't know exactly when it was going to happen. I wrote this code to profile our service every 20 seconds and if we had more goroutines than before, it will save the profile in samples folder.

### Run
```
git clone git@github.com:mhrlife/goroutineprofiler.git
cd goroutineprofiler
go run main.go
```
