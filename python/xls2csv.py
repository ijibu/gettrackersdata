#!/usr/bin/env python
#coding=utf-8
#writer:parming    Email：ming5536@163.com
#Python程序：批理转化Excel成CSV文件，依赖pyExcelerator模块
#$url http://www.cnblogs.com/ming5536/archive/2011/09/19/2181563.html
from pyExcelerator import *
import sys
import glob
class batxls2csv:
    def __init__(self):
                pass
    def savecsv1(self,arg):
        matrixgolb = []
        for sheet_name, values in parse_xls(arg, 'cp1251'): # parse_xls(arg) -- default encoding
            matrix = [[]]
            for row_idx, col_idx in sorted(values.keys()):
                #print row_idx,col_idx
                #print matrix
                v = values[(row_idx, col_idx)]
                if isinstance(v, unicode):
                    v = v.encode('cp866', 'backslashreplace')
                else:
                    v = str(v)
                last_row, last_col = len(matrix), len(matrix[-1])
                #下一行修改过
                while last_row <=row_idx:
                    matrix.extend([[]])
                    last_row = len(matrix)
                
                while last_col < col_idx:
                    matrix[-1].extend([''])
                    last_col = len(matrix[-1])
                
                matrix[-1].extend([v])
            for row in matrix:
                 csv_row = ','.join(row)
                 matrixgolb.append(csv_row)
        return matrixgolb        
        print  matrixgolb
    def savecsv2(self):
        
        filelist = glob.glob("*.xls")
        for filenam in filelist:
            matrixgolb=self.savecsv1(filenam)
            namecsv=filenam[:-4]+'.csv'
            file_object = open(namecsv, 'w+')
            for item in matrixgolb:
                file_object.write(item)
                file_object.write('\n')
                #print item
            file_object.close( )
        print 'ok!'
            
if __name__ == "__main__":
    test=batxls2csv()
    test.savecsv2()