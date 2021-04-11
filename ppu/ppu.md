PPU的作用：用来绘图的芯片

PPU的概念：
Background: 背景
Sprite: 精灵

PatternTable: 
存储的是图像数据。
总线映射：0x0000~0x1FFF, 8k
应该全部都映射到卡带上

VRAM的4k指的控制BG的4K

NameTable: 
控制的是背景图像。
总线映射：0x2000~0x2FFF, 4k
NameTable为的PatternTable的索引

AttributeTable：(64byte)
控制当前图像高两位palette的偏移量。


控制背景图像的1k显存：
960byte是NameTable, 64byte是Attribute Table
Attribute Table:
控制背景图像的颜色


PPU的内存映射

PPU的寄存器

图像布局：
像素:256x240 
一个tile: 8x8
所以整个图像的Tile分布: 32x30
0x2000~0x23FF: 第2kb,控制第一个图像窗格
0x2400~0x27FF: 第2kb,控制第二个图像窗格
0x2800~0x2bFF: 第3kb,
0x2C00~0x2FFF: 第4kb,

外加一个PPU滚动寄存器，就可以实现背景运动的效果了。

PPU只支持64种颜色，所以一种颜色只需要1byte即可。
作为对比，表示一种RGB颜色需要24byte,所以在颜色上压缩了3倍。

Palette(调色板)
背景调色版和精灵调色板分别16字节
0x3F00~0x3F0F: 背景调色板
0x3F10~0x3F1F: 精灵调色板

这些调色版每一个字节表示一种颜色(颜色的索引).

那么一个像素点只要用4bit索引调色板的颜色。而且这4bit还可以通过共用进一步压缩。

相比RGB的24byte, 压缩了48倍





