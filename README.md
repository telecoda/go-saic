#go-saic

Written by @telecoda

Image mosaic generator in go-lang

I know there are a gazillion free image mosaic utilies out there.  I decided to write this as an interesting coding exercise.

To get command parameters type:

    go-saic -?

##Command options

Work in progress.. this will change a lot.

    Usage of go-saic:
      -mosaic_image_path="image.png": path of image to create a mosaic from
      -r=false: search image directories recursively
      -s=false: search for images
      -source_dir="images": directory for source images
      -t=false: create thumbnails
      -target_image_path="target.png": path of mosaic image to be created
      -target_width=1024: default width of image to produce, height will be calculated to maintain aspect ratio
      -thumb_dir="thumbnail_images": directory to produce thumbnails images in
      -tile_height=32: default height of mosaic tile
      -tile_width=32: default width of mosaic tile
    

##Example usage

###search for pictures, create thumbnail and create a mosaic image

    go-saic -s -source_dir data/sourceimages -r -t -mosaic_image_path=data/testimage.png
    
This command will search the data/sourceimages directory recursively for images.

It will create thumbnails images of any images found and they will be saved in the "thumbnail_images" (Default) directory.

Image "data/testimage.png" will be used as a source image to be converted into an image mosaic.

###create a mosaic image


    go-saic -mosaic_image_path=data/testimage.png -tile_height=64 -tile_width=32 
    
    
tile_width + tile_height parameters are optional