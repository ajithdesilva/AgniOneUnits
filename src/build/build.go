/*
*****************************************************************************************************
# Author        :   Ajith Nilantha de Silva (ajithdesilva@zaion.ai) | 29/01/2024
# Copyright  	:	Zaion.AI 2024
# Class/module  :   build
# Objective     :   This package has ability to store the build information that fed buy the build time
#  					 using the -ldflags.
#######################################################################################################
# Author                        Date        Action      Description
#------------------------------------------------------------------------------------------------------
# Ajith de Silva               19/01/2024  Created     Created the initial version
#######################################################################################################
******************************************************************************************************
*/

package build

var Version string /// store the version

var Time string /// store the build time

var User string /// store the build user

var BuildGoVersion string /// store go version that used to build
