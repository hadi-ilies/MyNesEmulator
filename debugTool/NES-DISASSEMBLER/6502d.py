#!/usr/bin/python3

# 6502 NES disassembler Version May 10, 2018
# for assembly with asm6
# Doug Fraker 2017-2018

# Permission is hereby granted, free of charge, to any person obtaining a copy of 
# this software and associated documentation files (the "Software"), to deal in the 
# Software without restriction, including without limitation the rights to use, copy, 
# modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, 
# and to permit persons to whom the Software is furnished to do so, subject to the 
# following conditions:
#
# The above copyright notice and this permission notice shall be included in all 
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
# INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A 
# PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT 
# HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF  
# CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE 
# OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import sys
import os

global count #count within the small bank
global bankSize
global workArray
global workArraySmall
global currentBank

if len(sys.argv) < 2:
	print("usage: " + sys.argv[0] + " <path>")
	exit()
path = sys.argv[1]

# initialize some variables

count = 0
bankSize = 16384 # default size
first = ""
second = ""
third = ""
currentBank = ""



# define some functions

def ToASM(byte1,byte2,byte3):
	global count
	global currentBank
	count2 = 0
		
	if byte1 == "00":
		return ("\tbrk\t\t\t\t; 00") # none
		
	elif byte1 == "01":
		count += 1
		return ("\tora ($" + byte2 + ", x)\t; 01 " + byte2) # (Indirect,X)
		
	elif byte1 == "05":
		count += 1
		return ("\tora $" + byte2 + "\t\t\t; 05 " + byte2) # Zeropage
		
	elif byte1 == "06":
		count += 1
		return ("\tasl $" + byte2 + "\t\t\t; 06 " + byte2) # Zeropage
		
	elif byte1 == "08":
		return ("\tphp\t\t\t\t; 08 ") # none
		
	elif byte1 == "09":
		count += 1
		return ("\tora #$" + byte2 + "\t\t; 09 " + byte2)	# immediate
		
	elif byte1 == "0a":
		return ("\tasl a\t\t\t; 0a") # A

	elif byte1 == "0d":
		count += 2
		if (byte3 == "00"):
			return (".hex 0d "+byte2+" "+byte3)
		return ("\tora $" + byte3 + byte2 + "\t\t; 0d " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "0e":
		count += 2
		if (byte3 == "00"):
			return (".hex 0e "+byte2+" "+byte3)
		return ("\tasl $" + byte3 + byte2 + "\t\t; 0e " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "10":
		y = int(byte2, 16)
		if y > 127:
			y -= 256
		count2 = count + y + 2
		z = str(hex(count2))
		z = z[2:] 
		z = z.zfill(4)
		
		count += 1
		return ("\tbpl B" +currentBank+"_"+ z + " ; 10 " + byte2) # Relative
		
	elif byte1 == "11":
		count += 1
		 
		return ("\tora ($" + byte2 + "), y\t; 11 " + byte2) # (Indirect),Y
		
	elif byte1 == "15":
		count += 1
		 
		return ("\tora $" + byte2 + ", x\t\t; 15 " + byte2) # Zeropage, x
		
	elif byte1 == "16":
		count += 1
		 
		return ("\tasl $" + byte2 + ", x\t\t; 16 " + byte2) # Zeropage, x
		
	elif byte1 == "18":
		return ("\tclc\t\t\t\t; 18 ") # none
		
	elif byte1 == "19":
		count += 2
		if (byte3 == "00"):
			return (".hex 19 "+byte2+" "+byte3)
		return ("\tora $" + byte3 + byte2 + ", y\t; 19 " + byte2 + " " + byte3) # absolute, y
		
	elif byte1 == "1d":
		count += 2
		if (byte3 == "00"):
			return (".hex 1d "+byte2+" "+byte3)
		return ("\tora $" + byte3 + byte2 + ", x\t; 1d " + byte2 + " " + byte3) # absolute, x
		
	elif byte1 == "1e":
		count += 2
		if (byte3 == "00"):
			return (".hex 1e "+byte2+" "+byte3)
		return ("\tasl $" + byte3 + byte2 + ", x\t; 1e " + byte2 + " " + byte3) # absolute, x
		
	elif byte1 == "20":
		count += 2
		if (byte3 == "00"):
			return (".hex 20 "+byte2+" "+byte3)
		return ("\tjsr $" + byte3 + byte2 + "\t\t; 20 " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "21":
		count += 1
		 
		return ("\tand ($" + byte2 + ", x)\t; 21 " + byte2) # (Indirect,X)
		
	elif byte1 == "24":
		count += 1
		 
		return ("\tbit $" + byte2 + "\t\t\t; 24 " + byte2) # Zeropage
		
	elif byte1 == "25":
		count += 1
		 
		return ("\tand $" + byte2 + "\t\t\t; 25 " + byte2) # Zeropage
		
	elif byte1 == "26":
		count += 1
		 
		return ("\trol $" + byte2 + "\t\t\t; 26 " + byte2) # Zeropage
		
	elif byte1 == "28":
		return ("\tplp\t\t\t\t; 28 ") # none
		
	elif byte1 == "29":
		count += 1
		 
		return ("\tand #$" + byte2 + "\t\t; 29 " + byte2)	# immediate
		
	elif byte1 == "2a":
		return ("\trol a\t\t\t; 2a") # A
		
	elif byte1 == "2c":
		count += 2
		if (byte3 == "00"):
			return (".hex 2c "+byte2+" "+byte3)
		return ("\tbit $" + byte3 + byte2 + "\t\t; 2c " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "2d":
		count += 2
		if (byte3 == "00"):
			return (".hex 2d "+byte2+" "+byte3)
		return ("\tand $" + byte3 + byte2 + "\t\t; 2d " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "2e":
		count += 2
		if (byte3 == "00"):
			return (".hex 2e "+byte2+" "+byte3)
		return ("\trol $" + byte3 + byte2 + "\t\t; 2e " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "30":
		y = int(byte2, 16)
		if y > 127:
			y -= 256
		count2 = count + y + 2
		z = str(hex(count2))
		z = z[2:] 
		z = z.zfill(4)
		
		count += 1
		return ("\tbmi B" +currentBank+"_"+ z + " ; 30 " + byte2) # Relative
		
	elif byte1 == "31":
		count += 1
		 
		return ("\tand ($" + byte2 + "), y\t; 31 " + byte2) # (Indirect),Y
		
	elif byte1 == "35":
		count += 1
		 
		return ("\tand $" + byte2 + ", x\t\t; 35 " + byte2) # Zeropage, x
		
	elif byte1 == "36":
		count += 1
		 
		return ("\trol $" + byte2 + ", x\t\t; 36 " + byte2) # Zeropage, x
		
	elif byte1 == "38":
		return ("\tsec\t\t\t\t; 38 ") # none
		
	elif byte1 == "39":
		count += 2
		if (byte3 == "00"):
			return (".hex 39 "+byte2+" "+byte3)
		return ("\tand $" + byte3 + byte2 + ", y\t; 39 " + byte2 + " " + byte3) # absolute, y
		
	elif byte1 == "3d":
		count += 2
		if (byte3 == "00"):
			return (".hex 3d "+byte2+" "+byte3)
		return ("\tand $" + byte3 + byte2 + ", x\t; 3d " + byte2 + " " + byte3) # absolute, x
		
	elif byte1 == "3e":
		count += 2
		if (byte3 == "00"):
			return (".hex 3e "+byte2+" "+byte3)
		return ("\trol $" + byte3 + byte2 + ", x\t; 3e " + byte2 + " " + byte3) # absolute, x
		
	elif byte1 == "40":
		return ("\trti\t\t\t\t; 40 ") # none
		
	elif byte1 == "41":
		count += 1
		 
		return ("\teor ($" + byte2 + ", x)\t; 41 " + byte2) # (Indirect,X)
	
	elif byte1 == "45":
		count += 1
		 
		return ("\teor $" + byte2 + "\t\t\t; 45 " + byte2) # Zeropage
		
	elif byte1 == "46":
		count += 1
		 
		return ("\tlsr $" + byte2 + "\t\t\t; 46 " + byte2) # Zeropage
	
	elif byte1 == "48":
		return ("\tpha\t\t\t\t; 48 ") # none
		
	elif byte1 == "49":
		count += 1
		 
		return ("\teor #$" + byte2 + "\t\t; 49 " + byte2)	# immediate
		
	elif byte1 == "4a":
		return ("\tlsr a\t\t\t; 4a") # A
		
	elif byte1 == "4c":
		count += 2
		if (byte3 == "00"):
			return (".hex 4c "+byte2+" "+byte3)
		return ("\tjmp $" + byte3 + byte2 + "\t\t; 4c " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "4d":
		count += 2
		if (byte3 == "00"):
			return (".hex 4d "+byte2+" "+byte3)
		return ("\teor $" + byte3 + byte2 + "\t\t; 4d " + byte2 + " " + byte3) # absolute
	
	elif byte1 == "4e":
		count += 2
		if (byte3 == "00"):
			return (".hex 4e "+byte2+" "+byte3)
		return ("\tlsr $" + byte3 + byte2 + "\t\t; 4e " + byte2 + " " + byte3) # absolute
	
	elif byte1 == "50":
		y = int(byte2, 16)
		if y > 127:
			y -= 256
		count2 = count + y + 2
		z = str(hex(count2))
		z = z[2:] 
		z = z.zfill(4)
		
		count += 1
		return ("\tbvc B" +currentBank+"_"+ z + " ; 50 " + byte2) # Relative
		
		
	elif byte1 == "51":
		count += 1
		 
		return ("\teor ($" + byte2 + "), y\t; 51 " + byte2) # (Indirect),Y
		
	elif byte1 == "55":
		count += 1
		 
		return ("\teor $" + byte2 + ", x\t\t; 55 " + byte2) # Zeropage, x
		
	elif byte1 == "56":
		count += 1
		 
		return ("\tlsr $" + byte2 + ", x\t\t; 56 " + byte2) # Zeropage, x
	
	elif byte1 == "58":
		return ("\tcli\t\t\t\t; 58 ") # none
		
	elif byte1 == "59":
		count += 2
		if (byte3 == "00"):
			return (".hex 59 "+byte2+" "+byte3)
		return ("\teor $" + byte3 + byte2 + ", y\t; 59 " + byte2 + " " + byte3) # absolute, y
	
	elif byte1 == "5d":
		count += 2
		if (byte3 == "00"):
			return (".hex 5d "+byte2+" "+byte3)
		return ("\teor $" + byte3 + byte2 + ", x\t; 5d " + byte2 + " " + byte3) # absolute, x
		
	elif byte1 == "5e":
		count += 2
		if (byte3 == "00"):
			return (".hex 5e "+byte2+" "+byte3)
		return ("\tlsr $" + byte3 + byte2 + ", x\t; 5e " + byte2 + " " + byte3) # absolute, x
	
	elif byte1 == "60":
		return ("\trts\t\t\t\t; 60 ") # none
		
	elif byte1 == "61":
		count += 1
		 
		return ("\tadc ($" + byte2 + ", x)\t; 61 " + byte2) # (Indirect,X)
	
	elif byte1 == "65":
		count += 1
		 
		return ("\tadc $" + byte2 + "\t\t\t; 65 " + byte2) # Zeropage
		
	elif byte1 == "66":
		count += 1
		 
		return ("\tror $" + byte2 + "\t\t\t; 66 " + byte2) # Zeropage
	
	elif byte1 == "68":
		return ("\tpla\t\t\t\t; 68 ") # none

	elif byte1 == "69":
		count += 1
		 
		return ("\tadc #$" + byte2 + "\t\t; 69 " + byte2)	# immediate
		
	elif byte1 == "6a":
		return ("\tror a\t\t\t; 6a") # A
	
	elif byte1 == "6c":
		count += 2
		if (byte3 == "00"):
			return (".hex 6c "+byte2+" "+byte3)
		return ("\tjmp ($" + byte3 + byte2 + ")\t\t; 6c " + byte2 + " " + byte3) # absolute (indirect)
		
	elif byte1 == "6d":
		count += 2
		if (byte3 == "00"):
			return (".hex 6d "+byte2+" "+byte3)
		return ("\tadc $" + byte3 + byte2 + "\t\t; 6d " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "6e":
		count += 2
		if (byte3 == "00"):
			return (".hex 6e "+byte2+" "+byte3)
		return ("\tror $" + byte3 + byte2 + "\t\t; 6e " + byte2 + " " + byte3) # absolute
	
	elif byte1 == "70":
		y = int(byte2, 16)
		if y > 127:
			y -= 256
		count2 = count + y + 2
		z = str(hex(count2))
		z = z[2:] 
		z = z.zfill(4)
		
		count += 1
		return ("\tbvs B" +currentBank+"_"+ z + " ; 70 " + byte2) # Relative
		
	elif byte1 == "71":
		count += 1
		 
		return ("\tadc ($" + byte2 + "), y\t; 71 " + byte2) # (Indirect),Y
		
	elif byte1 == "75":
		count += 1
		 
		return ("\tadc $" + byte2 + ", x\t\t; 75 " + byte2) # Zeropage, x
		
	elif byte1 == "76":
		count += 1
		 
		return ("\tror $" + byte2 + ", x\t\t; 76 " + byte2) # Zeropage, x
	
	elif byte1 == "78":
		return ("\tsei\t\t\t\t; 78 ") # none
	
	elif byte1 == "79":
		count += 2
		if (byte3 == "00"):
			return (".hex 79 "+byte2+" "+byte3)
		return ("\tadc $" + byte3 + byte2 + ", y\t; 79 " + byte2 + " " + byte3) # absolute, y
		
	elif byte1 == "7d":
		count += 2
		if (byte3 == "00"):
			return (".hex 7d "+byte2+" "+byte3)
		return ("\tadc $" + byte3 + byte2 + ", x\t; 7d " + byte2 + " " + byte3) # absolute, x
		
	elif byte1 == "7e":
		count += 2
		if (byte3 == "00"):
			return (".hex 7e "+byte2+" "+byte3)
		return ("\tror $" + byte3 + byte2 + ", x\t; 7e " + byte2 + " " + byte3) # absolute, x
	
	elif byte1 == "81":
		count += 1
		 
		return ("\tsta ($" + byte2 + ", x)\t; 81 " + byte2) # (Indirect,X)
		
	elif byte1 == "84":
		count += 1
		 
		return ("\tsty $" + byte2 + "\t\t\t; 84 " + byte2) # Zeropage
	
	elif byte1 == "85":
		count += 1
		 
		return ("\tsta $" + byte2 + "\t\t\t; 85 " + byte2) # Zeropage
		
	elif byte1 == "86":
		count += 1
		 
		return ("\tstx $" + byte2 + "\t\t\t; 86 " + byte2) # Zeropage
		
	elif byte1 == "88":
		return ("\tdey\t\t\t\t; 88 ") # none
		
	elif byte1 == "8a":
		return ("\ttxa\t\t\t\t; 8a ") # none
		
	elif byte1 == "8c":
		count += 2
		if (byte3 == "00"):
			return (".hex 8c "+byte2+" "+byte3)
		return ("\tsty $" + byte3 + byte2 + "\t\t; 8c " + byte2 + " " + byte3) # absolute	
		
	elif byte1 == "8d":
		count += 2
		if (byte3 == "00"):
			return (".hex 8d "+byte2+" "+byte3)
		return ("\tsta $" + byte3 + byte2 + "\t\t; 8d " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "8e":
		count += 2
		if (byte3 == "00"):
			return (".hex 8e "+byte2+" "+byte3)
		return ("\tstx $" + byte3 + byte2 + "\t\t; 8e " + byte2 + " " + byte3) # absolute	
		
	elif byte1 == "90":
		y = int(byte2, 16)
		if y > 127:
			y -= 256
		count2 = count + y + 2
		z = str(hex(count2))
		z = z[2:] 
		z = z.zfill(4)
		
		count += 1
		return ("\tbcc B" +currentBank+"_"+ z + " ; 90 " + byte2) # Relative
		
	elif byte1 == "91":
		count += 1
		 
		return ("\tsta ($" + byte2 + "), y\t; 91 " + byte2) # (Indirect),Y
		
	elif byte1 == "94":
		count += 1
		 
		return ("\tsty $" + byte2 + ", x\t\t; 94 " + byte2) # Zeropage, x
	
	elif byte1 == "95":
		count += 1
		 
		return ("\tsta $" + byte2 + ", x\t\t; 95 " + byte2) # Zeropage, x
		
	elif byte1 == "96":
		count += 1
		 
		return ("\tstx $" + byte2 + ", y\t\t; 96 " + byte2) # Zeropage, y
		
	elif byte1 == "98":
		return ("\ttya\t\t\t\t; 98 ") # none
		
	elif byte1 == "99":
		count += 2
		if (byte3 == "00"):
			return (".hex 99 "+byte2+" "+byte3)
		return ("\tsta $" + byte3 + byte2 + ", y\t; 99 " + byte2 + " " + byte3) # absolute, y
		
	elif byte1 == "9a":
		return ("\ttxs\t\t\t\t; 9a ") # none
		
	elif byte1 == "9d":
		count += 2
		if (byte3 == "00"):
			return (".hex 9d "+byte2+" "+byte3)
		return ("\tsta $" + byte3 + byte2 + ", x\t; 9d " + byte2 + " " + byte3) # absolute, x
		
	elif byte1 == "a0":
		count += 1
		 
		return ("\tldy #$" + byte2 + "\t\t; a0 " + byte2) # immediate
		
	elif byte1 == "a1":
		count += 1
		 
		return ("\tlda ($" + byte2 + ", x)\t; a1 " + byte2) # (Indirect,X)
		
	elif byte1 == "a2":
		count += 1
		 
		return ("\tldx #$" + byte2 + "\t\t; a2 " + byte2) # immediate
		
	elif byte1 == "a4":
		count += 1
		 
		return ("\tldy $" + byte2 + "\t\t\t; a4 " + byte2) # Zeropage
		
	elif byte1 == "a5":
		count += 1
		 
		return ("\tlda $" + byte2 + "\t\t\t; a5 " + byte2) # Zeropage
		
	elif byte1 == "a6":
		count += 1
		 
		return ("\tldx $" + byte2 + "\t\t\t; a6 " + byte2) # Zeropage
		
	elif byte1 == "a8":
		return ("\ttay\t\t\t\t; a8 ") # none	
		
	elif byte1 == "a9":
		count += 1
		 
		return ("\tlda #$" + byte2 + "\t\t; a9 " + byte2) # immediate
		
	elif byte1 == "aa":
		return ("\ttax\t\t\t\t; aa ") # none
		
	elif byte1 == "ac":
		count += 2
		if (byte3 == "00"):
			return (".hex ac "+byte2+" "+byte3)
		return ("\tldy $" + byte3 + byte2 + "\t\t; ac " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "ad":
		count += 2
		if (byte3 == "00"):
			return (".hex ad "+byte2+" "+byte3)
		return ("\tlda $" + byte3 + byte2 + "\t\t; ad " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "ae":
		count += 2
		if (byte3 == "00"):
			return (".hex ae "+byte2+" "+byte3)
		return ("\tldx $" + byte3 + byte2 + "\t\t; ae " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "b0":
		y = int(byte2, 16)
		if y > 127:
			y -= 256
		count2 = count + y + 2
		z = str(hex(count2))
		z = z[2:] 
		z = z.zfill(4)
		
		count += 1
		return ("\tbcs B" +currentBank+"_"+ z + " ; b0 " + byte2) # Relative
		
	elif byte1 == "b1":
		count += 1
		 
		return ("\tlda ($" + byte2 + "), y\t; b1 " + byte2) # (Indirect),Y
		
	elif byte1 == "b4":
		count += 1
		 
		return ("\tldy $" + byte2 + ", x\t\t; b4 " + byte2) # Zeropage, x
		
	elif byte1 == "b5":
		count += 1
		 
		return ("\tlda $" + byte2 + ", x\t\t; b5 " + byte2) # Zeropage, x
		
	elif byte1 == "b6":
		count += 1
		 
		return ("\tldx $" + byte2 + ", y\t\t; b6 " + byte2) # Zeropage, y
		
	elif byte1 == "b8":
		return ("\tclv\t\t\t\t; b8 ") # none
		
	elif byte1 == "b9":
		count += 2
		if (byte3 == "00"):
			return (".hex b9 "+byte2+" "+byte3)
		return ("\tlda $" + byte3 + byte2 + ", y\t; b9 " + byte2 + " " + byte3) # absolute, y
		
	elif byte1 == "ba":
		return ("\ttsx\t\t\t\t; ba ") # none
		
	elif byte1 == "bc":
		count += 2
		if (byte3 == "00"):
			return (".hex bc "+byte2+" "+byte3)
		return ("\tldy $" + byte3 + byte2 + ", x\t; bc " + byte2 + " " + byte3) # absolute, x
		
	elif byte1 == "bd":
		count += 2
		if (byte3 == "00"):
			return (".hex bd "+byte2+" "+byte3)
		return ("\tlda $" + byte3 + byte2 + ", x\t; bd " + byte2 + " " + byte3) # absolute, x
		
	elif byte1 == "be":
		count += 2
		if (byte3 == "00"):
			return (".hex be "+byte2+" "+byte3)
		return ("\tldx $" + byte3 + byte2 + ", y\t; be " + byte2 + " " + byte3) # absolute, y
		
	elif byte1 == "c0":
		count += 1
		 
		return ("\tcpy #$" + byte2 + "\t\t; c0 " + byte2) # immediate
		
	elif byte1 == "c1":
		count += 1
		 
		return ("\tcmp ($" + byte2 + ", x)\t; c1 " + byte2) # (Indirect,X)
		
	elif byte1 == "c4":
		count += 1
		 
		return ("\tcpy $" + byte2 + "\t\t\t; c4 " + byte2) # Zeropage	
		
	elif byte1 == "c5":
		count += 1
		 
		return ("\tcmp $" + byte2 + "\t\t\t; c5 " + byte2) # Zeropage
	
	elif byte1 == "c6":
		count += 1
		 
		return ("\tdec $" + byte2 + "\t\t\t; c6 " + byte2) # Zeropage
		
	elif byte1 == "c8":
		return ("\tiny\t\t\t\t; c8 ") # none
		
	elif byte1 == "c9":
		count += 1
		 
		return ("\tcmp #$" + byte2 + "\t\t; c9 " + byte2) # immediate
		
	elif byte1 == "ca":
		return ("\tdex\t\t\t\t; ca ") # none
		
	elif byte1 == "cc":
		count += 2
		if (byte3 == "00"):
			return (".hex cc "+byte2+" "+byte3)
		return ("\tcpy $" + byte3 + byte2 + "\t\t; cc " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "cd":
		count += 2
		if (byte3 == "00"):
			return (".hex cd "+byte2+" "+byte3)
		return ("\tcmp $" + byte3 + byte2 + "\t\t; cd " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "ce":
		count += 2
		if (byte3 == "00"):
			return (".hex ce "+byte2+" "+byte3)
		return ("\tdec $" + byte3 + byte2 + "\t\t; ce " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "d0":
		y = int(byte2, 16)
		if y > 127:
			y -= 256
		count2 = count + y + 2
		z = str(hex(count2))
		z = z[2:] 
		z = z.zfill(4)
		
		count += 1
		return ("\tbne B" +currentBank+"_"+ z + " ; d0 " + byte2) # Relative	
		
	elif byte1 == "d1":
		count += 1
		 
		return ("\tcmp ($" + byte2 + "), y\t; d1 " + byte2) # (Indirect),Y
		
	elif byte1 == "d5":
		count += 1
		 
		return ("\tcmp $" + byte2 + ", x\t\t; d5 " + byte2) # Zeropage, x
		
	elif byte1 == "d6":
		count += 1
		 
		return ("\tdec $" + byte2 + ", x\t\t; d6 " + byte2) # Zeropage, x
		
	elif byte1 == "d8":
		return ("\tcld\t\t\t\t; b8 ") # none
		
	elif byte1 == "d9":
		count += 2
		if (byte3 == "00"):
			return (".hex d9 "+byte2+" "+byte3)
		return ("\tcmp $" + byte3 + byte2 + ", y\t; d9 " + byte2 + " " + byte3) # absolute, y
		
	elif byte1 == "dd":
		count += 2
		if (byte3 == "00"):
			return (".hex dd "+byte2+" "+byte3)
		return ("\tcmp $" + byte3 + byte2 + ", x\t; dd " + byte2 + " " + byte3) # absolute, x
		
	elif byte1 == "de":
		count += 2
		if (byte3 == "00"):
			return (".hex de "+byte2+" "+byte3)
		return ("\tdec $" + byte3 + byte2 + ", x\t; de " + byte2 + " " + byte3) # absolute, x
		
	elif byte1 == "e0":
		count += 1
		 
		return ("\tcpx #$" + byte2 + "\t\t; e0 " + byte2) # immediate
	
	elif byte1 == "e1":
		count += 1
		 
		return ("\tsbc ($" + byte2 + ", x)\t; e1 " + byte2) # (Indirect,X)
	
	elif byte1 == "e4":
		count += 1
		 
		return ("\tcpx $" + byte2 + "\t\t\t; e4 " + byte2) # Zeropage
		
	elif byte1 == "e5":
		count += 1
		 
		return ("\tsbc $" + byte2 + "\t\t\t; e5 " + byte2) # Zeropage
	
	elif byte1 == "e6":
		count += 1
		 
		return ("\tinc $" + byte2 + "\t\t\t; e6 " + byte2) # Zeropage
		
	elif byte1 == "e8":
		return ("\tinx\t\t\t\t; e8 ") # none
		
	elif byte1 == "e9":
		count += 1
		 
		return ("\tsbc #$" + byte2 + "\t\t; e9 " + byte2)	# immediate
	
	elif byte1 == "ea":
		return ("\tnop\t\t\t\t; ea ") # none
		
	elif byte1 == "ec":
		count += 2
		if (byte3 == "00"):
			return (".hex ec "+byte2+" "+byte3)
		return ("\tcpx $" + byte3 + byte2 + "\t\t; ec " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "ed":
		count += 2
		if (byte3 == "00"):
			return (".hex ed "+byte2+" "+byte3)
		return ("\tsbc $" + byte3 + byte2 + "\t\t; ed" + byte2 + " " + byte3) # absolute
	
	elif byte1 == "ee":
		count += 2
		if (byte3 == "00"):
			return (".hex ee "+byte2+" "+byte3)
		return ("\tinc $" + byte3 + byte2 + "\t\t; ee " + byte2 + " " + byte3) # absolute
		
	elif byte1 == "f0":
		y = int(byte2, 16)
		if y > 127:
			y -= 256
		count2 = count + y + 2
		z = str(hex(count2))
		z = z[2:] 
		z = z.zfill(4)
		
		count += 1
		return ("\tbeq B" +currentBank+"_"+ z + " ; f0 " + byte2) # Relative	
		
	elif byte1 == "f1":
		count += 1
		 
		return ("\tsbc ($" + byte2 + "), y\t; f1 " + byte2) # (Indirect),Y
	
	elif byte1 == "f5":
		count += 1
		 
		return ("\tsbc $" + byte2 + ", x\t\t; f5 " + byte2) # Zeropage, x	
		
	elif byte1 == "f6":
		count += 1
		 
		return ("\tinc $" + byte2 + ", x\t\t; f6 " + byte2) # Zeropage, x
		
	elif byte1 == "f8":
		return ("\tsed\t\t\t\t; f8 ") # none
		
	elif byte1 == "f9":
		count += 2
		if (byte3 == "00"):
			return (".hex f9 "+byte2+" "+byte3)
		return ("\tsbc $" + byte3 + byte2 + ", y\t; f9 " + byte2 + " " + byte3) # absolute, y
		
	elif byte1 == "fd":
		count += 2
		if (byte3 == "00"):
			return (".hex fd "+byte2+" "+byte3)
		return ("\tsbc $" + byte3 + byte2 + ", x\t; fd " + byte2 + " " + byte3) # absolute, x
	
	elif byte1 == "fe":
		count += 2
		if (byte3 == "00"):
			return (".hex fe "+byte2+" "+byte3)
		return ("\tinc $" + byte3 + byte2 + ", x\t; fe " + byte2 + " " + byte3) # absolute, x	
		
	else:
		return (".db $" + byte1) # unknown opcode

#end of big function






# START OF PROGRAM

filename = os.path.basename(path)

try:
	fileIn = open(path, "rb") #read bytes
except:
	print("\nERROR: couldn't find file\n")
	raise
	
print (filename)
filesize = os.path.getsize(path)
print("filesize = ", filesize)
folder = os.path.dirname(path)
	
workArray = fileIn.read() #make a big int array 

testarray = bytearray(b'\x4e\x45\x53\x1a') # NES 1a



# validate header	
	
for i in range (0,4):
	if (workArray[i] != testarray[i]):
		print("\nERROR: couldn't find iNES header\n")
		exit()
		
		
# get ROM sizes
	
prgROM = workArray[4]
prgROMtotal = prgROM * 0x4000
print ("PRGROM size = ", prgROM, " = ", prgROMtotal)

chrROM = workArray[5]
chrROMtotal = chrROM * 0x2000
print ("CHRROM size = ", chrROM, " = ", chrROMtotal)	
	
a = 16 + prgROMtotal + chrROMtotal
print ("Header + PRGROM + CHRROM = ", a)
if (filesize != a):		
	print ("\nERROR: filesize does not match the header")
	if (filesize < a):
		exit()
	else:
		print ("Will try to disassemble anyway.\n") 
else:
	print ("filesize matches header, ok")
	
	
# get mapper	

a = 0
b = 0
c = 0
byte6 = 0
byte7 = 0

byte6 = workArray[6]
a = byte6 >> 4
byte7 = workArray[7]
b = byte7 & 0xf0
c = a + b

Map = ""

# default bankSize = 16384, already set
if (prgROM == 2):
	bankSize = 32768
#1/2 = 8192

if c == 0:
	Map = "NROM"
elif c == 1:
	Map = "MMC1 SxROM"
elif c == 2:
	Map = "UxROM"
elif c == 3:
	Map = "CNROM"
elif c == 4:
	Map = "MMC3 TxROM"	
	bankSize = 8192
elif c == 5:
	Map = "MMC5 ExROM"
elif c == 7:
	Map = "AxROM"
	bankSize = 32768
elif c == 9:
	Map = "MMC2 PxROM"
elif c == 10:
	Map = "MMC4 FxROM"
elif c == 11:
	Map = "COLOR DREAMS"
elif c == 13:
	Map = "CPROM"
elif c == 16:
	Map = "Bandai"
elif c == 18:
	Map = "Jaleco"
elif c == 19:
	Map = "Namco 163"
elif c == 20:
	Map = "FDS"
elif c == 21:
	Map = "Konami VRC4"
elif c == 22:
	Map = "Konami VRC2"	
elif c == 23:
	Map = "Konami variation"
elif c == 24:
	Map = "Konami VRC6"
elif c == 25:
	Map = "Konami variation"
elif c == 26:
	Map = "Konami VRC6"
elif c == 28:
	Map = "Action 53"
elif c == 30:
	Map = "UNROM RetroUSB"
elif c == 31:
	Map = "NSF music"
elif c == 32:
	Map = "Irem's G-101"
elif c == 33:
	Map = "Taito's TC0190"
elif c == 34:
	Map = "BNROM or NINA-001"
	bankSize = 32768
elif c == 36:
	Map = "TXC"
elif c == 48:
	Map = "Taito's TC0690"
elif c == 64:
	Map = "Tengen RAMBO-1"
elif c == 65:
	Map = "Irem's H3001"
elif c == 66:
	Map = "GxROM or MHROM"
elif c == 67:
	Map = "Sunsoft-3"
elif c == 68:
	Map = "Sunsoft-4"
elif c == 69:
	Map = "Sunsoft FME-7"
elif c == 70:
	Map = "Bandai"
elif c == 71:
	Map = "Codemasters"
elif c == 72:
	Map = "Jaleco's JF-17"
elif c == 73:
	Map = "Konami VRC3"
elif c == 74:
	Map = "Eastern games"
elif c == 75:
	Map = "Konami VRC1"
elif c == 76:
	Map = "Namcot 108"
elif c == 77:
	Map = "Irem"
elif c == 78:
	Map = "Irem"
elif c == 79:
	Map = "NINA-03 or NINA-06"
elif c == 80:
	Map = "Taito's X1-005"	
elif c == 82:
	Map = "Taito's X1-017"	
elif c == 85:
	Map = "Konami VRC7"	
elif c == 86:
	Map = "Jaleco's JF-13"	
elif c == 87:
	Map = "Jaleco"	
elif c == 88:
	Map = "Namco"		
elif c == 89:
	Map = "Sunsoft"	
elif c == 93:
	Map = "Sunsoft"	
elif c == 94:
	Map = "HVC-UN1ROM"	
elif c == 99:
	Map = "Vs. System"	
elif c == 118:
	Map = "TKSROM and TLSROM"	
elif c == 119:
	Map = "TQROM"		
	
else:
	Map = "Other / Too Lazy to type them all in."


print ("Mapper number = ", c, " = ", Map)	


# mirroring = low bit of byte6

a = byte6 & 0x08 # 4 screen
if (a == 0):
	a = byte6 & 0x01 # 2 screen
	
if (a == 0):
	print ("horizontal mirroring")
elif (a == 1):
	print ("vertical mirroring")
else:
	print ("4 screen mode")
	
	
# extra RAM at 6000 = byte6 ? bit 2

a = byte6 & 0x02
if (a != 0):
	print ("extra RAM at $6000, yes")
	
	
# sanity check	

if prgROM == 0 or filesize < 16400:
	print ("file too small, not valid")
	exit()
	
	
	
# split ROM into 2 binary files, PRG minus the header (called .bin), and CHR

newName = os.path.splitext(filename)[0] # strip the extension
newPath = os.path.join(folder, newName + ".bin")	
	
fileOut = open(newPath,"wb") # write bytes

for i in range (16, prgROMtotal+16):
	a = workArray[i]
	a = bytes([a])
	fileOut.write(a)
	
fileOut.close	
print (newName+ ".bin created")




#chrROMtotal
if (chrROM != 0):
	#newName = os.path.splitext(filename)[0] # strip the extension
	newPath = os.path.join(folder, newName + ".chr")	
	
	fileOut = open(newPath,"wb") # write bytes
	
	for i in range (prgROMtotal+16, chrROMtotal+prgROMtotal+16):
		a = workArray[i]
		a = bytes([a])
		fileOut.write(a)
	
	fileOut.close	
	
	print (newName+ ".chr created")
	
else:
	print ("No CHR")

	
# get bank size, from user

Valid = 0
b = 0
print("Recommended bank size = ", bankSize)
while (Valid == 0):
	a = input("OK? Y/N:")
	if a == "Y" or a == "y":
		Valid = 1
	else:
		while (b == 0):
			b = input("1 = 8192, 2 = 16384, 4 = 32768:")
			if b == "1":
				bankSize = 8192
				Valid = 1
				break
			elif b == "2":
				bankSize = 16384
				Valid = 1
				break
			elif b == "4":
				bankSize = 32768
				Valid = 1
				break	
			else:
				b = 0

				
if (bankSize > prgROMtotal):
	print("exceeds total PRG ROM size...")
	bankSize = prgROMtotal
	
print("bankSize = ", bankSize)



# start writing the MAIN ASM file

#newName = os.path.splitext(filename)[0] # strip the extension
newPath = os.path.join(folder, newName + ".asm")	
	
fileOutMain = open(newPath,"w") # write text
print (newName+ ".asm created")

	
fileOutMain.write ("; " + filename + " disassembly\n")
fileOutMain.write ("; for asm6\n\n")

fileOutMain.write ("; *** HEADER ***\n\n")
fileOutMain.write (".db \"NES\", $1a\n")


a = workArray[4] # byte 4
c = str(a)
fileOutMain.write (".db " + c + " ; = number of PRG banks * $4000\n")

a = workArray[5] # byte 5
c = str(a)
fileOutMain.write (".db " + c + " ; = number of CHR banks * $2000\n")

a = workArray[6] # byte 6
c = str(a)
fileOutMain.write (".db " + c + "\t; " + Map + "\n")

a = workArray[7] # byte 7
c = str(a)
fileOutMain.write (".db " + c + "\n")

a = workArray[8] # byte 8
c = str(a)
fileOutMain.write (".db " + c + "\n")

a = workArray[9] # byte 9
c = str(a)
fileOutMain.write (".db " + c + "\n")

a = workArray[10] # byte 10
c = str(a)
fileOutMain.write (".db " + c + "\n")
fileOutMain.write (".db 0,0,0,0,0\n\n") # bytes 11-15


fileOutMain.write ("; *** PRG ROM ***\n\n")

if prgROM > 1:
	fileOutMain.write (".base $8000\n\n")	# default starting address
else:
	fileOutMain.write (".base $c000\n\n")	# default starting address

workArraySmall = [0] * bankSize



# start writing the other ASM files, bank by bank

bankNumberTotal = int (prgROMtotal / bankSize)

for bankNumber in range (0,bankNumberTotal):
	currentBank = str(bankNumber)
	fileOutMain.write(".include "+newName + currentBank + ".asm\n\n")
	newPath = os.path.join(folder, newName + currentBank + ".asm")	
	
	fileOutSmall = open(newPath,"w+") # write text, and read it
	print (newName+currentBank+ ".asm created")

	fileOutSmall.write (";"+newName+currentBank+"\n\n\n\n")
	
	#create a smaller array
	for i in range (0,bankSize):
		j = i + 16 + (bankNumber*bankSize)
		workArraySmall[i] = workArray[j] # note both int arrays
	
	# decode the array
	count = 0
	while (count < bankSize-2): # change later ?
		a = workArraySmall[count]	# get 3 bytes, just in case
		first = str (hex (a)) #convert int to hex string
		first = first[2:] # strip the 0x off
		first = first.zfill(2) # at least 2 wide, fill zero
		a = workArraySmall[count+1]
		second = str (hex (a))
		second = second[2:]
		second = second.zfill(2)
		a = workArraySmall[count+2]
		third = str (hex (a))
		third = third[2:]
		third = third.zfill(2)
		
		z = str(hex(count))
		z = z[2:] 
		z = z.zfill(4)
		
		fileOutSmall.write("B"+currentBank+"_"+z+":\t")
		
		outString = ToASM(first,second,third)
		fileOutSmall.write(outString+"\n")
		
		count += 1
	
	# print the final bytes... if needed
	if (count < bankSize):
		a = workArraySmall[count]	# get 3 bytes, just in case
		first = str (hex (a)) #convert int to hex string
		first = first[2:] # strip the 0x off
		first = first.zfill(2) # at least 2 wide, fill zero
		fileOutSmall.write("\t\t.db $" + first+"\n")
		count += 1
	if (count < bankSize):
		a = workArraySmall[count]	# get 3 bytes, just in case
		first = str (hex (a)) #convert int to hex string
		first = first[2:] # strip the 0x off
		first = first.zfill(2) # at least 2 wide, fill zero
		fileOutSmall.write("\t\t.db $" + first+"\n")
		count += 1	
		
	
	
	# remove broken labels	
	
	fileOutSmall.seek(0) # needed ?
	contents = fileOutSmall.readlines()
	listAll = []
	for i in range(len(contents)):
		listAll.extend(contents[i].split())
		
	listLabels = [] # make a list of all labels in sub-file
	
	word = ""
	last = ""
	position = 0
	
	loop = len(listAll)
	for i in range (0,loop):
		word = str(listAll[i])
		last = word[-1:]
		if last == ":":
			word = word [:-1]
			listLabels.append(word)
	
	# see if reference to label, if not, remove it.
	fileOutSmall.seek(0)
	filedata = fileOutSmall.read()
	
	for i in range (0,loop):
		
		word = str(listAll[i])
		if word == "bcc" or word == "bcs" or word == "bvc" or word == "bvs" or word == "beq" or word == "bne" or word == "bmi" or word == "bpl":
			word2 = str(listAll[i+1])
			if word2 not in listLabels:

				# kill the word in the the original text file now
				fullword = word+" "+word2+" ;"
				filedata = filedata.replace(fullword, ";removed\n\t.hex ") # replace it with this
				

			
	fileOutSmall.seek(0)
	fileOutSmall.write(filedata)
	fileOutSmall.close
	
	if bankNumberTotal > bankNumber+1:
		fileOutMain.write (".base $8000\n\n")	# default starting address, maybe fix this later

	# end of sub bank asm decode loop	


	
fileOutMain.write ("\n\n; *** CHR ROM ***\n\n")	
if (chrROM != 0):
	fileOutMain.write (".incbin "+newName+".chr\n\n")
else:
	fileOutMain.write (";No CHR ROM\n\n")
fileOutMain.close
fileIn.close					

print ("done!")
