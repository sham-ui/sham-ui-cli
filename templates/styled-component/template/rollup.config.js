import { babel } from '@rollup/plugin-babel';
import shamUICompiler from 'rollup-plugin-sham-ui-templates';
import copy from 'rollup-plugin-copy';
import scss from 'rollup-plugin-scss';
import pkg from './package.json';

export default [ {
    input: './src/{{name}}.sfc',
    external: [
        'sham-ui',
        'sham-ui-macro/ref.macro',
        /@babel\/runtime/,
        /@corejs/
    ],
    output: [
        { file: pkg.main, format: 'cjs', exports: 'auto' },
        { file: pkg.module, format: 'es' }
    ],
    plugins: [
        shamUICompiler( {
            extensions: [ '.sfc' ],
            compilerOptions: {
                asModule: false,
                asSingleFileComponent: true,
                removeDataTest: false
            }
        } ),
        babel( {
            extensions: [ '.js', '.sht', '.sfc' ],
            exclude: [ 'node_modules/**' ],
            babelHelpers: 'inline'
        } )
    ]
}, {
    input: 'src/styles/{{name}}.scss',
    plugins: [
        copy( {
            targets: [
                { src: 'src/styles/{{name}}.scss', dest: './', rename: 'style.scss' }
            ]
        } ),
        scss( {
            failOnError: true,
            output: 'style.css',
            watch: 'src/styles',
            sass: require( 'sass' )
        } )
    ]
} ];
