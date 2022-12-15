const multimatch = require( 'multimatch' );
const async = require( 'async' );
const render = require( 'consolidate' ).handlebars.render;

function renameFiles( skipInterpolation ) {
    if ( typeof skipInterpolation === 'string' ) {
        skipInterpolation = [ skipInterpolation ];
    }

    return function( files, metalsmith, done ) {
        const keys = Object.keys( files );
        const metalsmithMetadata = metalsmith.metadata();

        async.each( keys, function( file, next ) {
            if ( skipInterpolation &&
                multimatch( [ file ], skipInterpolation, { dot: true } ).length ) {
                return next();
            }

            if ( !/{{([^{}]+)}}/g.test( file ) ) {
                return next();
            }

            render( file, metalsmithMetadata, function( err, res ) {
                if ( err ) {
                    err.message = `[${file}] ${err.message}`;
                    return next( err );
                }
                const content = files[ file ].contents.toString();
                delete files[ file ];
                files[ res ] = {
                    contents: Buffer.from( content )
                };
                next();
            } );

        }, done );
    }
}

module.exports = renameFiles;
