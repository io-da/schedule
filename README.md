# [Go](https://golang.org/) Schedule
A schedule utility to generate all the dates.  

[![Build Status](https://travis-ci.org/io-da/schedule.svg?branch=master)](https://travis-ci.org/io-da/schedule)
[![Maintainability](https://api.codeclimate.com/v1/badges/3bf3737ea61c79b5d74a/maintainability)](https://codeclimate.com/github/io-da/schedule/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/3bf3737ea61c79b5d74a/test_coverage)](https://codeclimate.com/github/io-da/schedule/test_coverage)
[![GoDoc](https://godoc.org/github.com/io-da/schedule?status.svg)](https://godoc.org/github.com/io-da/schedule)

## Installation
``` go get github.com/io-da/schedule ```

## Overview
1. [Schedule](#Schedule)
2. [CronExpression](#CronExpression)
3. [IteratorExpression](#IteratorExpression)
4. [CronInstance](#CronInstance)  
5. [Putting it together](#Putting-it-together)

## Introduction
This library is intended to be used as an utility for generating _time.Time_ structs.    
It should provide an alternative to _Linux_'s _Cron_ and _Crontab_ with some extra spice.  
**This library purposely does not implement the actual job handling/execution.**  
Clean codebase. **No reflection, no closures.**

## Getting Started

### Schedule
_Schedule_ is the struct used to represent a set of retrievable _time.Time_ structs.  
Optionally it can also contain a CronExpression that will only start producing _time.Time_ structs after the scheduled ones.  
Schedules can be initialized in 2 ways: ```schedule.At(at ...time.Time)``` or ```schedule.In(in ...time.Duration)```.  
```schedule.At``` creates a _Schedule_ that returns the provided _time.Time_ structs.  
```schedule.In``` creates a _Schedule_ that returns ```time.Now()``` plus the provided _time.Duration_ values.  
The _time.Time_ structs can be retrieved by executing the function ```sch.Following()```.  
Moving the _Schedule_ state forward is achieved by executing the function ```sch.Next()```.

Optionally it is also possible to provide a _CronExpression_ to a _Schedule_ (```sch.Cron(crn *CronExpression)```).  
The _CronExpression_ will only start to be used after the schedules times.
  
### CronExpression
_CronExpression_ struct represents a full crontab expression.  
An expression is initialized using the function ```schedule.Cron()```.  
This expression is defined by using a set of functions designed for readability.  
The _CronExpression_ does not generate times it self. For that we need a _CronInstance_.  

##### Examples
_At 00:00 on day-of-month 29 and on Sunday in February_. (0 0 29 2 0):
```go
Cron().
OnMonths(time.February).
OnDays(29).
OnWeekdays(time.Sunday)
```

### IteratorExpression
Iterator expressions are any type that implements the _IteratorExpression_ interface.  
These are used to represent sets of values.
```go
type IteratorExpression interface {
    Next(from int, inc bool) (int, bool)
    Contains(val int) bool
}
```
The library provides 2 different ones.  
_BetweenExpression_:
```go
type BetweenExpression struct {
    x    int
    y    int
    step int
}
```
This expression also allows configuration of it's stepping value using ```Every(s int)```.  

_ListExpression_:
```go
type ListExpression struct {
    values []int
}
```

##### Examples
(Between) _Every 2 hours_: ```Cron().OnHours(BetweenHours(0,22).Every(2))```  
(List) _Specifically at 0 and 12 hours_: ```Cron().OnHours(ListHours(0,12))```  

### CronInstance
To start generating _time.Time_ structs, we just need to create a _CronInstance_ from the _CronExpression_ ```crn.NewInstance(from time.Time)```.  
It is required to provide a from _time.Time_ for the CronInstance to be able to identify its following _time.Time_.  
The _CronInstance_ exposes only 2 functions.
```go
Next() error
Following() time.Time
```

```Next()``` is used to determine the next following _time.Time_. Every execution advances it's internal state.  
```Following()``` is used to retrieve the determined _time.Time_. It can be retrieved without any consequence. 

## Putting it together
```go
import (
    "github.com/io-da/schedule"
)

func main() {
    // instantiate a schedule that produces a time.Time one hour in the future
    sch := schedule.In(time.Hour)
    
    // setup a cron to start producing time.Time every hour thereafter
    sch.Cron(schedule.Cron().EveryHour())
    
    for {
        // check if the schedule/cron expired
        if err := sch.Next(); err == nil {
            // and potentially handle the error
            hypotheticalLogger.error(err)
            break
        }

        // retrieve the determined time.Time
        at := sch.Following()
        
        // and handle it
        hypotheticalTaskManager.Add(hypotheticalTask, at)
    }
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)