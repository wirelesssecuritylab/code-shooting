#!/usr/bin/env python
import sys
SPRING = 1
SUMMER = 2
AUTUMN = 3
WINTER = 4


class Season:
    def __init__(self, season: int):
        self._season = season

    def __str__(self):
        if self._season == SPRING:
            return 'Spring'
        elif self._season == SUMMER:
            return 'Summer'
        elif self._season == AUTUMN:
            return 'Autumn'
        elif self._season == WINTER:
            return 'Winter'
        else:
            raise Exception('invalid season')
