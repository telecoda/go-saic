#go-saic

Written by @telecoda

Image mosaic generator in go-lang

I know there are a gazillion free image mosaic utilies out there.  I decided to write this as an interesting coding exercise.

To get command parameters type:

    go-saic -?
    
##Building
First fetch all the dependencies:


    go get -u -v
   
    go build

That's it.   

##Command options

    Usage of ./go-saic:
      -c=false: Analyse thumbnail images for most prominent colour
      -d=false: search for images in source_dir
      -f="image.png": path of input image (used to create mosaic from)
      -m=false: Create a photo mosaic image
      -o="output.png": path of output image
      -output_width=1024: default width of image to produce, height will be calculated to maintain aspect ratio
      -r=false: search image directories recursively
      -source_dir="data/input/sourceimages": directory for source images
      -t=false: create thumbnails
      -thumb_dir="data/output/thumbnail_images": directory to create thumbnails images in
      -tile_height=32: height of image tiles in output image
      -tile_width=32: width of image tiles in output image
      
##Example usage

###search for pictures, create thumbnail and create a mosaic image

    go-saic -d -source_dir data/input/sourceimages/ -r -t -f=data/input/testimages/testimage.png -m    
This command will discover images in the data/input/sourceimages directory recursively.

It will create thumbnails images of any images found and they will be saved in the "./data/output/thumbnail_images" (Default) directory.

Image "data/testimage.png" will be used as a source image to be converted into an image mosaic.

###create a mosaic image

	go-saic -m -f=data/input/testimages/testimage.png -tile_height=64 -tile_width=32 
    
    
tile_width + tile_height parameters are optional



##How it works?
This section describes the basic functionality of the software:


###Step one: Source image discovery (-d)

* Prerequisites:- needs a directory containing images.

This step must be performed at least once.  This is used to discover and catalog the images that will be
available as a reference that can be used to create a mosaic from.

###Step two: Source image thumbnail creation (-t)
* Prerequisites:- needs "discovery" to have run to produce a list of images to process.

This step must be performed at least once.  It must always be run after step one (discovery).
The process creates smaller scaled thumbnails of the source images in a separate working directory.

Steps one and two can be run in isolation if this is a long running task.

###Step three: Thumbnail image colour analysis (-c)
* Prerequisites:- needs "thumbnailing" to have run to produce a list of images to process.

This step will analyse the thumbails that have been produced and create and calculate the most prominent colour in the image.

###Step four: Creation of a photo mosaic (-m)
* Prerequisites:- needs "discovery" & "thumbnail" to have run.
			    need input image - this is the image that will be used a the basis of the photo mosaic

This step will create a photo mosaic using a source image.

The process will not update the source "mosaic_image" a new "target_image" will be created.

Summary of the photo mosaic process:

* target image is created as a copy of the source mosaic_image (this can be scaled to a different size)
* divide target image into a number of "tiles" based upon the tile height and width parameters
* each tile is analysed to find its prominent colour
* each tile is replaced with a thumbnail image of a similar colour
* repeat for all the tiles on the image
* probably have lots of gaps in resulting image due to lack of photos
* think of a crafty way of filling the gaps...    

