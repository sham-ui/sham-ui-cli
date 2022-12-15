const util = require( 'util' );
const fs = require( 'fs' );
const path = require( 'path' );
const { exec } = require( 'child_process' );
const glob = require( 'glob-fs' );
const findup = require( 'findup-sync' );
const ignore = require( 'parse-gitignore' );
const mm = require( 'micromatch' );
const inquirer = require( 'inquirer' );
const { gray } = require( 'chalk' );

const copyFilePromise = util.promisify( fs.copyFile );
const unlinkFilePromise = util.promisify( fs.unlink );
const execPromise = util.promisify( exec );
const mkdirPromise = util.promisify( fs.mkdir );

function eachSeries( arr, iteratorFn ) {
    return arr.reduce(
        ( acc, item ) => acc.then(
            () => iteratorFn( item )
        ),
        Promise.resolve()
    );
}

function filterSeries( arr, fn ) {
    return arr.reduce(
        ( acc, item ) => acc.then(
            ( accRes ) => fn( item ).then( value => {
                if ( value ) {
                    accRes.push( item );
                }
                return accRes;
            } )
        ),
        Promise.resolve( [] )
    );
}

function parseGitignore( file ) {
    const opts = this.setDefaults( this.pattern.options, {} );

    const gitignoreFile = findup(
        '.gitignore',
        {
            cwd: path.dirname( file.toAbsolute() )
        }
    );
    const ignorePatterns = ignore( gitignoreFile );

    if ( mm.any( file.relative, ignorePatterns, opts ) ) {
        file.exclude = true;
    }
    return file;
}

function getFiles( dir ) {
    return glob( { builtins: false } )
        .use( parseGitignore )
        .readdirSync( '/**/*.*', {
            cwd: dir
        } );
}


function copyFiles( srcDir, destDir, files ) {
    return Promise.all(
        files.map(
            x => mkdirPromise(
                path.dirname( path.join( destDir, x ) ),
                { recursive: true }
            ).then( () => copyFilePromise(
                path.join( srcDir, x ),
                path.join( destDir, x )
            ) )
        )
    );
}

function openInEditor( editor, file ) {
    return execPromise( `${ editor } ${ file }` ).then(
        () => {},
        ( e ) => console.error( e )
    )
}

function removeFiles( editor, dirFile, files ) {
    return eachSeries(
        files,
        x => inquirer.prompt( [ {
            name: 'answer',
            type: 'expand',
            message: `Remove file ${gray( x )}`,
            default: 'remove',
            expanded: true,
            choices: [
                {
                    key: 'y',
                    name: 'Remove',
                    value: 'remove',
                },
                {
                    key: 'e',
                    name: 'Open in editor',
                    value: 'open',
                },
                {
                    key: 's',
                    name: 'Skip',
                    value: 'skip',
                }
            ]
        } ] ).then( ( { answer } ) => {
            const filePath = path.join( dirFile, x );
            if ( 'remove' === answer ) {
                return unlinkFilePromise( filePath )
            }
            if ( 'open' === answer ) {
                return openInEditor( editor, filePath );
            }
        } )
    );
}

function intersection( tmpFiles, destFiles ) {
    return destFiles.filter( x => tmpFiles.includes( x ) );
}

function readFile( f ) {
    return new Promise( ( resolve, reject ) => {
        fs.readFile( f, function( err, data ) {
            if ( err ) {
                reject( err );
            } else {
                resolve( data );
            }
        } )
    } );
}

function overrideFiles( editor, srcDir, destDir, files ) {
    return filterSeries(
        files,
        x => Promise.all( [
            readFile( path.join( srcDir, x ) ),
            readFile( path.join( destDir, x ) )
        ] ).then(
            ( [ src, dest ] ) => !src.equals( dest )
        )
    ).then(
        filteredFiles => eachSeries(
            filteredFiles,
            x => inquirer.prompt( [ {
                name: 'answer',
                type: 'expand',
                message: `Override file ${ gray( x )}`,
                default: 'open',
                expanded: true,
                choices: [
                    {
                        key: 'y',
                        name: 'Override',
                        value: 'override',
                    },
                    {
                        key: 'e',
                        name: 'Open in editor',
                        value: 'open',
                    },
                    {
                        key: 's',
                        name: 'Skip',
                        value: 'skip',
                    },
                ]
            } ] ).then( ( { answer } ) => {
                if ( 'override' === answer ) {
                    return copyFilePromise(
                        path.join( srcDir, x ),
                        path.join( destDir, x )
                    )
                }
                if ( 'open' === answer ) {
                    return Promise.all( [
                        openInEditor( editor, path.join( destDir, x ) ),
                        openInEditor( editor, path.join( srcDir, x ) )
                    ] );
                }
            } )
        )
    );
}

function onlyInSecond( first, second ) {
    return second.filter( x => !first.includes( x ) );
}

function upgrade( tmpPath, destPath, done ) {
    inquirer.prompt( [ {
        name: 'editor',
        type: 'input',
        message: `Editor for open files`
    } ] ).then( ( { editor } ) => (
        {
            editor,
            tmpFiles: getFiles( tmpPath ),
            destFiles: getFiles( destPath )
        }
    ) ).then( ctx => {
        const forCopy = onlyInSecond( ctx.destFiles, ctx.tmpFiles );
        return copyFiles( tmpPath, destPath, forCopy )
            .then(
                () => ctx
            );
    } ).then( ctx => {
        const forRemove = onlyInSecond( ctx.tmpFiles, ctx.destFiles );
        return removeFiles( ctx.editor, destPath, forRemove )
            .then(
                () => ctx
            );
    } ).then( ctx => {
        const forOverride = intersection( ctx.tmpFiles, ctx.destFiles );
        return overrideFiles( ctx.editor, tmpPath, destPath, forOverride ).then( () => null );
    } ).then(
        done,
        done
    );
}


module.exports = upgrade;
