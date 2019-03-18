package image

import (
    "strings"

    "golang.org/x/net/context"
    "github.com/sirupsen/logrus"
    "github.com/docker/docker/client"
    "github.com/seveirbian/comparator/compare"
    "github.com/docker/docker/daemon/graphdriver/overlay2"
)

type ImageComparator struct {
    Image1Name string
    Image1Tag string
    Image2Name string
    Image2Tag string
    Dir1ID string
    Dir2ID string

    Ctx context.Context
    Client *client.Client

    Comparator *compare.Comparator
}

func Init(image1 string, image2 string) (*ImageComparator, error) {
    name1Slices := strings.Split(image1, ":")
    if len(name1Slices) != 2 {
        logrus.Warn("Image name is not valid...")
    }

    name2Slices := strings.Split(image2, ":")
    if len(name2Slices) != 2 {
        logrus.Warn("Image name is not valid...")
    }

    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
    if err != nil {
        logrus.Warn("Fail to create docker client...")
        return nil, err
    }

    image1Info, _, err := cli.ImageInspectWithRaw(ctx, image1)
    if err != nil {
        logrus.Warnf("Fail to inspect image: %s\n", image1)
        return nil, err
    }

    image2Info, _, err := cli.ImageInspectWithRaw(ctx, image2)
    if err != nil {
        logrus.Warnf("Fail to inspect image: %s\n", image2)
        return nil, err
    }

    var dir1ID string
    dir1ID = image1Info.GraphDriver.Data["UpperDir"]
    dir1ID = strings.Split(dir1ID, "/var/lib/docker/overlay2/")[1]
    dir1ID = strings.Split(dir1ID, "/diff")[0]

    var dir2ID string
    dir2ID = image2Info.GraphDriver.Data["UpperDir"]
    dir2ID = strings.Split(dir2ID, "/var/lib/docker/overlay2/")[1]
    dir2ID = strings.Split(dir2ID, "/diff")[0]

    return &ImageComparator {
        Image1Name: name1Slices[0], 
        Image1Tag: name1Slices[1], 
        Image2Name: name2Slices[0], 
        Image2Tag: name2Slices[1], 
        Dir1ID: dir1ID, 
        Dir2ID: dir2ID, 

        Ctx:ctx, 
        Client: cli, 
    }, nil
}

func (i *ImageComparator) Compare() error {
    driver, err := overlay2.Init("/var/lib/docker/overlay2", []string{}, nil, nil)
    if err != nil {
        logrus.WithField("err", err).Warn("Fail to create overlay2 driver...")
        return err
    }

    image1MountPath, err := driver.Get(i.Dir1ID, "")
    if err != nil {
        logrus.WithField("err", err).Warn("Fail to mount overlayfs...")
        return err
    }
    defer driver.Put(i.Dir1ID)

    image1Path := image1MountPath.Path()

    image2MountPath, err := driver.Get(i.Dir2ID, "")
    if err != nil {
        logrus.WithField("err", err).Warn("Fail to mount overlayfs...")
        return err
    }
    defer driver.Put(i.Dir2ID)

    image2Path := image2MountPath.Path()

    comparator, err := compare.Init(image1Path, image2Path)
    if err != nil {
        logrus.WithField("err", err).Warn("Fail to init comparator...")
        return err
    }

    err = comparator.Compare()
    if err != nil {
        logrus.WithField("err", err).Warn("Fail to compare...")
        return err
    }

    i.Comparator = comparator

    return nil
}


















