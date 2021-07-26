const async = require( 'async' );
const render = require( 'consolidate' ).handlebars.render;

function renderTemplates( skipInterpolation ) {
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

            const str = files[ file ].contents.toString();

            if ( !/{{([^{}]+)}}/g.test( str ) ) {
                return next();
            }

            render( str, metalsmithMetadata, function( err, res ) {
                if ( err ) {
                    err.message = `[${file}] ${err.message}`;
                    return next( err );
                }
                files[ file ].contents = Buffer.from( res );
                next();
            } );

        }, done );
    }
}

module.exports = renderTemplates;
