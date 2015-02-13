# github-issues-dater [![Build Status](https://travis-ci.org/nabeken/github-issues-dater.png?branch=master)](https://travis-ci.org/nabeken/github-issues-dater)

Github Issues has search syntax.
For example, we can query issues that is updated within a week like this:
`is:issue is:open updated:2015-02-07..2015-02-14` (Today is 2015-02-14)

You must build a date range for updated manually.

github-issues-dater allows you to build `updated` filed with a relative date.

## URL

Use `http://github-issues-dater.herokuapp.com` instead of `https://github.com`.

For example, if you want to query issues updated within a week:

https://github.com/GoogleCloudPlatform/kubernetes/issues?q=is%3Aissue+is%3Aopen+updated%3A2015-02-07..2015-02-14

will be

http://github-issues-dater.herokuapp.com/GoogleCloudPlatform/kubernetes/issues?q=is%3Aissue+is%3Aopen+updated%3A1w
