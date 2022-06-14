import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import renderer from 'sham-ui-test-helpers';
import { onFocusIn } from '../../../../../../src/directives/on-focus-in';
// eslint-disable-next-line max-len
import RoutesArticlesCommonArticleForm from '../../../../../../src/components/routes/articles/common/article-form.sfc';

it( 'renders correctly', async() => {
    const DI = createDI();
    DI.bind( 'store', {
        articleCategories: jest.fn().mockReturnValueOnce(
            Promise.resolve( {
                categories: [
                    { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                    { 'id': '6', 'name': 'Название категории', 'slug': 'nazvanie-kategorii' }
                ],
                meta: {}
            } )
        ),
        articleTags: jest.fn().mockReturnValueOnce(
            Promise.resolve( {
                tags: [
                    { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                    { 'id': '6', 'name': 'Название тега', 'slug': 'nazvanie-tega' }
                ],
                meta: {}
            } )
        )
    } );
    const meta = renderer( RoutesArticlesCommonArticleForm, {
        DI,
        directives: {
            onFocusIn,
            ...directives
        },
        publishedAt: new Date( '2022-04-29T12:24:51.637Z' )
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );


it( 'fail load data', async() => {
    const DI = createDI();
    DI.bind( 'store', {
        articleCategories: jest.fn().mockReturnValueOnce(
            Promise.resolve( {
                categories: [
                    { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                    { 'id': '6', 'name': 'Название категории', 'slug': 'nazvanie-kategorii' }
                ],
                meta: {}
            } )
        ),
        articleTags: jest.fn().mockReturnValueOnce(
            Promise.reject( {} )
        )
    } );
    const meta = renderer( RoutesArticlesCommonArticleForm, {
        DI,
        directives: {
            onFocusIn,
            ...directives
        },
        publishedAt: new Date( '2022-04-29T12:24:51.637Z' )
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
