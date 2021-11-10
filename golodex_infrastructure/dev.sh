#!/bin/bash
export golodex_dir=$(pwd)
cd ~
rm -f golodex
ln -s $golodex_dir golodex