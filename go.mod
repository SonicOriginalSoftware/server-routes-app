module git.sonicoriginal.software/routes/app

go 1.19

require git.sonicoriginal.software/server v0.0.0

replace (
	git.sonicoriginal.software/routes/app => github.com/SonicOriginalSoftware/server-routes-app v0.0.0
	git.sonicoriginal.software/server => github.com/SonicOriginalSoftware/server v0.0.0
)
