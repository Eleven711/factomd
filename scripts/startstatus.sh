#!/usr/bin/env bash
tail -f out.txt | gawk -f scripts/status.awk