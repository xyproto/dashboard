# Dashboard [![Build Status](https://travis-ci.com/xyproto/dashboard.svg?branch=master)](https://travis-ci.com/xyproto/dashboard) [![GoDoc](https://godoc.org/github.com/xyproto/dashboard?status.svg)](http://godoc.org/github.com/xyproto/dashboard)

![Desktop](https://raw.githubusercontent.com/xyproto/dashboard/master/screenshots/desktop.png)

Simple dashboard web application skeleton.

Using
-----

* [permissions2](https://github.com/xyproto/permissions2), for cookies and authentication, my own project
* [martini](https://github.com/go-martini/martini), for handling requests
* [render](https://github.com/unrolled/render), for rendering templates

And:

* [pure.css](http://purecss.io/), by Yahoo, for the layout and menu
* [Google Charts](https://developers.google.com/chart/), for the charts

This works reasonably well, but in the future negroni or gin will probably be chosen over martini, boostrap over pure and a more lightweight library than Google Charts for the charts. Other than that, the current choices are viable.

TODO
----

* Look at https://github.com/oal/admin for inspiration

General information
-------------------

* Version: 0.1
* License: MIT
* Alexander F. Rødseth
