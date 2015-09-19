#!/bin/bash
cd templates/static
esc -pkg templates -o ../templates.go .
cd - >/dev/null
