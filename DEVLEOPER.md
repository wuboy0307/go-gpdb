Table of Contents
=================

   * [Introduction](#introduction)
   * [PreChanges](#prechanges)
   * [Setup](#setup)
   * [Pull Request](#pull-request)

# Introduction

If you want to contribute or want to hack here are the details that can allow you to setup the environment to test and play.

# PreChanges.

Ensure you have opened a issue request before making any pull request to this repository

# Setup

+ Git clone the repository
    ```
    https://github.com/pivotal-gss/go-gpdb.git
    ```
+ Create a new branch
    ```
    git checkout -b <new-branch>
    ```
+ Setup your GOPATH & GOROOT environment variables  
+ Identify which cli, do you want to modify or hack like is it [datalab cli](https://github.com/pivotal-gss/go-gpdb/tree/master/datalab) or [gpdb cli](https://github.com/pivotal-gss/go-gpdb/tree/master/gpdb)
+ The instruction to setup each CLI is provided here
    + [Datalab cli](https://github.com/pivotal-gss/go-gpdb/tree/master/datalab#developers--contributors)
    + [Gpdb cli](https://github.com/pivotal-gss/go-gpdb/tree/master/gpdb#developers--contributors)
+ Its always best to test all the changes on a VM, so that the local machine is not having a lot of unnecessary stuff, checkout the [steps](https://github.com/pivotal-gss/go-gpdb/tree/master/test#setup) on how to provision a developer mode VM

# Pull Request

If you wish to provide the changes and make it part of this repository, please feel free to open a [pull request](https://github.com/pivotal-gss/go-gpdb/pulls)

**Happy Coding**