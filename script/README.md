## Run

download

```bash
# 5.9  - python 3.8  命令行安装失败执行此步骤
https://www.paraview.org/download/ 
```

```bash
# pvpython grid2paraview.py args ...
pvpython grid2paraview.py circle_grid.txt  circle_grid
```


## Error

You are using pip version 8.1.2, however version 23.3.2 is available. You should consider upgrading via the 'pip install --upgrade pip' command.

```bash
yum install python3-pip
pip3 install --upgrade pip
yum remove python-pip
sudo python3 -m pip install --upgrade pip
```

# 升级pip:
```bash
sudo wget https://bootstrap.pypa.io/pip/2.7/get-pip.py
sudo python get-pip.py
pip -V
```

# 升级pip3:
```bash
sudo wget https://bootstrap.pypa.io/pip/3.8/get-pip.py
sudo python3 get-pip.py
pip -V
```



Syntax: stl2surf.py stlfile surffile: 
```bash
python stl2surf.py apollo.stl apollo.surf
```

ImportError: No module named parallel_bucket_sort
```bash
使用绝对路径
```

File "log2txt.py", line 14, in <module>
path = os.environ["SPARTA_PYTHON_TOOLS"]
File "/usr/lib64/python2.7/UserDict.py", line 23, in __getitem__
raise KeyError(key)
```bash
export SPARTA_PYTHON_TOOLS=/usr/local/openmpi
```

Traceback (most recent call last):
  File "log2txt.py", line 16, in <module>
    from olog import olog
ImportError: No module named olog
```bash
pip install olog
```

ERROR: Could not find a version that satisfies the requirement olog (from versions: none)
ERROR: No matching distribution found for olog


