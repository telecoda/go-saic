#go-saic
<a href="https://dl.dropboxusercontent.com/u/13846060/go-saic-gopher.png"><img src="https://dl.dropboxusercontent.com/u/13846060/go-saic-gopher-scaled.png" alt="go-saic gopher" align="right"/></a>

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

    -R=false: repair & compact db
    -X=false: clear image db
    -d=false: search for images in source_dir
    -f="image.png": path of input image (used to create mosaic from)
    -l=false: list image db content
    -m=false: Create a photo mosaic image
    -o="output.png": path of output image
    -output_width=1024: default width of image to produce, height will be calculated to maintain aspect ratio
    -r=false: search image directories recursively
    -source_dir="data/input/sourceimages": directory for source images
    -t=false: create thumbnails
    -thumb_dir="data/output/thumbnail_images": directory to create thumbnails images in
    -tile_size=32: size of image tiles in output image, width & height are the same
    -type="matched": Type of mosaic (tinted or matched)      
##Example usage

###search for pictures, create thumbnail and create a mosaic image

    go-saic -d -source_dir data/input/sourceimages/ -r -t -f=data/input/testimages/testimage.png -m -type=tinted
This command will discover images in the data/input/sourceimages directory recursively.

It will create thumbnails images of any images found and they will be saved in the "./data/output/thumbnail_images" (Default) directory.

Image "data/testimage.png" will be used as a source image to be converted into an image mosaic.

The default output image name is "output.png"

###create a mosaic image

	go-saic -m -f=data/input/testimages/testimage.png -tile_size=64 -type=tinted
    
    
tile_size parameter is optional



##How it works?
This section describes the basic functionality of the software:


###Step one: Source image discovery (-d)

* Prerequisites:- needs a directory containing images.

This step must be performed at least once.  This is used to discover and catalog the images that will be available as a reference that can be used to create a mosaic from.

###Step two: Source image thumbnail creation (-t)
* Prerequisites:- needs "discovery" to have run to produce a list of images to process.

This step must be performed at least once.  It must always be run after step one (discovery).
The process creates smaller scaled thumbnails of the source images in a separate working directory.

Steps one and two can be run in isolation if this is a long running task.

###Step three: Thumbnail image colour analysis
This is run as part of the thumbnailing process.

This step will analyse the thumbails that have been produced and calculate the most prominent colour in the image.  Thumbnail details will be stored in the image DB.

###Step four: Creation of a photo mosaic (-m)
* Prerequisites:- needs "discovery" & "thumbnail" to have run.
			    need input image - this is the image that will be used a the basis of the photo mosaic

This step will create a photo mosaic using a source image.

The process will not update the "source_image" a new image will be created. (Option -o allows you to specific a filename)

go-saic supports two different methods for rendering mosaics

**"matched" & "tinted"**

##Image DB
go-saic creates a simple database of images during the discovery and thumbnailing processes.  This is stored locally in a JSON database using [tiedot nosql db in golang
](https://github.com/HouzuoGuo/tiedot)
If the database doesn't exist it will be created.  To trash and rebuild the database use:-

    go-saic -X
    
You can list the content of the db like this:-

    go-saic -l 


##Mosaic Types
go-saic supports two different photo mosaic types

The mosaic type is specified using the -type parameter
    
    -type=tinted
    -type=matched

###type=matched

To create a matched mosaic the following process occurs:-

* target image is created as a copy of the source mosaic_image (this can be scaled to a different size)
* divide target image into a number of "tiles" based upon the tile_size parameter
* each tile is analysed to find its prominent colour
* search thumbnails DB for images of a similar colour (this looks for a close match then gets progressively more relaxed.  Therefore if there are no decent matches you could end up with anything!)
* For each image tile, scale the thumbnail image to match the tile_size and draw it

###type=tinted
When using a small set of images you'll probably not get a decent match so a little jiggery-pokery is necessary....  This is how tinted mosaics came about.

To create a tinted mosaic the following process occurs:-

* target image is created as a copy of the source mosaic_image (this can be scaled to a different size)
* divide target image into a number of "tiles" based upon the tile_size parameter
* each tile is analysed to find its prominent colour
* get a list of ALL thumbnail images (we don't want anyone missed out..)
* for each image tile allocate a thumbnail image. If we run out of images, go back to start and repeat images (they'll never notice...)
* For each image tile, scale the thumbnail image to match the tile_size
* Then convert the scaled image to a greyscale versions
* Then create a separate transparent image with the prominent colour from the original image in this tile location
* Merge the two together and you get a nice tinted image of the scaled thumbnail
* It should bear some resemblance the original image

##Other notes
There are other options that let you tweak the output such as:

    -output_width=2048   // this will scale the size of the output image
    -tile_size=64        // this alters the size of tiles draw

##Thumbnail tips
Thumbnails are always square.  Maybe it was me being lazy but it makes the whole process MUCH easier.  The minimum dimension is chosen as the size and then a square is cropped from the centre of the image.

Thumbnails are always 128x128 pixel.  I was lazy and didn't want to add another parameter.  Maybe I will one day...

Finally:-  I NEVER delete thumbnails.  Using the -X option will delete the database but NOT the thumbnail images.  If they are not referenced in the DB they won't be used.

The reason I did this is that photos are very precious and I didn't want someone using my software with the wrong directory path at deleting their photos.  I'll leave it in YOUR capable hands to delete files you don't want....
