const path = require( 'path' );
const Metalsmith = require( 'metalsmith' );
const getOptions = require( './options/generate' );

const filterFiles = require( './metalsmith-plugins/filter-files' );
const promptQuestions = require( './metalsmith-plugins/prompt' );
const renameFiles = require( './metalsmith-plugins/rename-files' );
const removeFiles = require( './metalsmith-plugins/remove-files' );

const classify = require( './utils/classify' );

function destroy( name, src, dest, done ) {
    const opts = getOptions( name, src );
    const metalsmith = Metalsmith( path.join( src, 'template' ) );


    const inTestRelativePathChunks = [ '..', '..', 'src', 'components', `${name}${opts.testImportExtension}` ];
    name.split( '/' ).forEach(
        () => inTestRelativePathChunks.unshift( '..' )
    );

    const data = Object.assign( metalsmith.metadata(), {
        inPlace: true,
        noEscape: true,
        classifiedName: classify( name ),
        dest: dest,
        testRelativePath: inTestRelativePathChunks.join( '/' )
    } );

    metalsmith.use( promptQuestions( opts.prompts ) )
        .use( filterFiles( opts.filters ) )
        .use( renameFiles( opts.skipInterpolation ) )
        .use( removeFiles );

    metalsmith.clean( false )
        .source( '.' ) // start from template root instead of `./src` which is Metalsmith's default for `source`
        .destination( dest )
        .build( done );

    return data;
}

module.exports = destroy;
