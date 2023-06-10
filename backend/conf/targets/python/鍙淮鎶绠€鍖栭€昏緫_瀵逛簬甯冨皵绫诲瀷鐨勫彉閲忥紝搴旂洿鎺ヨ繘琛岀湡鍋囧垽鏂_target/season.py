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
        int_to_season = {SPRING: 'Spring', SUMMER: 'Summer', AUTUMN: 'Autumn', WINTER: 'Winter'}
        return int_to_season[self._season]

    @staticmethod
    def is_valid_season(season):
        if season not in [SPRING, SUMMER, AUTUMN, WINTER]:
            return False
        return True


if __name__ == "__main__":
    season_seq = sys.argv[1]
    is_valid_season = Season.is_valid_season(season_seq)
    if is_valid_season == False:
        print('invalid season')
