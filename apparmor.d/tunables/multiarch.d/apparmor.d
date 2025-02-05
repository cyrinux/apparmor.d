# apparmor.d - Full set of apparmor profiles
# Extended system directories definition
# Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# To allow extended personalisation without breaking everything.
# All apparmor profiles should always use the variables defined here.

# Universally unique identifier
@{uuid}=[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*

# Hexadecimal
@{hex}=[0-9a-fA-F]*

# Date and time
@{date}=[0-9][0-9][0-9][0-9]-[1-12]-[1-31]
@{time}=[1-24]-[0-60]-[0-60]

# @{MOUNTDIRS} is a space-separated list of where user mount directories
# are stored, for programs that must enumerate all mount directories on a
# system.
@{MOUNTDIRS}=/media/ @{run}/media/ /mnt/

# @{MOUNTS} is a space-separated list of all user mounted directories.
@{MOUNTS}=@{MOUNTDIRS}/*/

# Common places for binaries and libraries across distributions
@{bin}=/{usr/,}{s,}bin
@{lib}=/{usr/,}lib{,exec,32,64}
