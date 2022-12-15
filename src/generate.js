const path = require( 'path' );
const Metalsmith = require( 'metalsmith' );
const getOptions = require( './options/generate' );

const filterFiles = require( './metalsmith-plugins/filter-files' );
const promptQuestions = require( './metalsmith-plugins/prompt' );
const renderTemplates = require( './metalsmith-plugins/render-templates' );
const renameFiles = require( './metalsmith-plugins/rename-files' );

const classify = require( './utils/classify' );

function scaffold( name, src, dest, done ) {
    const opts = getOptions( name, src );
    const metalsmith = Metalsmith( path.join( src, 'template' ) );


    const inTestRelativePathChunks = [ '..', '..', 'src', 'components', `${name}${opts.testImportExtension}` ];
    const srcComponentRelativePath = [];
    name.split( '/' ).forEach(
        () => {
            inTestRelativePathChunks.unshift( '..' );
            srcComponentRelativePath.unshift( '..' );
        }
    );
    srcComponentRelativePath.pop();

    const data = Object.assign( metalsmith.metadata(), {
        inPlace: true,
        noEscape: true,
        classifiedName: classify( name ),
        testRelativePath: inTestRelativePathChunks.join( '/' ),
        srcComponentRelativePath: srcComponentRelativePath.join( '/' )
    } );

    metalsmith.use( promptQuestions( opts.prompts ) )
        .use( filterFiles( opts.filters ) )
        .use( renameFiles( opts.skipInterpolation ) )
        .use( renderTemplates( opts.skipInterpolation ) );

    metalsmith.clean( false )
        .source( '.' ) // start from template root instead of `./src` which is Metalsmith's default for `source`
        .destination( dest )
        .build( done );

    return data;
}

module.exports = scaffold;
