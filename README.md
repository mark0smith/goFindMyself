# goFindMyself
For Golang learning and memory practices.

## Usage
```
> ./hello -h
Usage of ./hello:
  -hint int
        If set, when recall test failes, hint will be given.
        0 for no hint, 1 for diff hint, 2 for full hint
  -m int
        Generated number wont't be bigger than this number. (default 100)
  -n int
        Number of Random Numbers. (default 30)
  -r rememberLogfile
        If set, generated numbers will be logged into rememberLogfile.
        You should set this if you want to do recall test later!
  -r_file string
        Filename of remember log (default "log.txt")
  -recall
        If set, run a recall test, instead of generating random numbers.
  -recallLog recallLogfile
        If set, recall info will be logged into recallLogfile. (default true)
  -recall_file string
        Filename of recall log (default "recall_log.txt")
  -u    If set, all generated numbers will be unique. (default true)
```

If running without any parameters, this program generates 30 random numbers in range of (0,100).
```
> ./hello
2023-10-05 22:06:51 Random Number is [48 17 29 83 97 96 14 44 78 74 35 90 58 38 68 33 77 63 85 2 3 25 81 84 9 32 61 5 53 37]
```

You can specific numbers count by `-n` and range maxium by `-m`.
```
> ./hello -n 4 -m 10
2023-10-05 22:08:28 Random Number is [3 0 6 1]
```

## Memory Practice

If you want to log these numbers for later recalling test, you should specific `-r`, thus numbers will be written into `log.txt`.
You can Specific filename by `-r_file`.
```
> ./hello -r
2023-10-05 22:10:58 Random Number is [64 57 8 54 79 15 97 39 37 7 13 11 67 6 78 80 36 60 32 83 70 59 25 27 92 35 9 75 1 22]
> tail -n1 log.txt
2023-10-05 22:10:58 [64 57 8 54 79 15 97 39 37 7 13 11 67 6 78 80 36 60 32 83 70 59 25 27 92 35 9 75 1 22]
```

If you believe you have remember these numbers and wanna have a recall test, you can add `-recall`.
If `-hint int` is set, some hint will be given if you are wrong, 0 for no hint, 1 for diff hint, 2 for full hint

```
> ./hello -recall                                                                                                                                                   
What do you remember?                                                                                                                                                                                             
13 40 70 56 96 55 82 63 92 71 85 83 46 1 4 7 67 87 61 57 81 18 72 58 31 36 52 53 78                                                                                                                               
                                                                                                                                                                                                                  
You have entered: 13 40 70 56 96 55 82 63 92 71 85 83 46 1 4 7 67 87 61 57 81 18 72 58 31 36 52 53 78                                                                                                             
Are you sure you remember it right?                                                                                                                                                                               
> ./hello -recall -hint 1
What do you remember?
13 40 70 56 96 55 82 63 92 71 85 83 46 1 4 7 67 87 61 57 81 18 72 58 31 36 52 53 78

You have entered: 13 40 70 56 96 55 82 63 92 71 85 83 46 1 4 7 67 87 61 57 81 18 72 58 31 36 52 53 78
Are you sure you remember it right?

Hint Part:
You are missing these numbers: 26
You add these numbers which should't exist: 
> ./hello -recall -hint 2
What do you remember?
13 40 70 56 96 55 82 63 92 71 85 83 46 1 4 7 67 87 61 57 81 18 72 58 31 36 52 53 78

You have entered: 13 40 70 56 96 55 82 63 92 71 85 83 46 1 4 7 67 87 61 57 81 18 72 58 31 36 52 53 78
Are you sure you remember it right?

Hint Part:
The Right: 13 40 70 56 96 55 82 63 92 71 26 85 83 46 1 4 7 67 87 61 57 81 18 72 58 31 36 52 53 78
The Wrong: 13 40 70 56 96 55 82 63 92 71 85 83 46 1 4 7 67 87 61 57 81 18 72 58 31 36 52 53 78
```


