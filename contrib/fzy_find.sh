#!/bin/sh
# use fzy(1) to find->open password

cpass open $(cpass ls | fzy)
