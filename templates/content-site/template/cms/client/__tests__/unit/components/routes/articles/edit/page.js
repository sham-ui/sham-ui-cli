import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import { onFocusIn } from '../../../../../../src/directives/on-focus-in';
// eslint-disable-next-line max-len
import RoutesArticlesEditPage  from '../../../../../../src/components/routes/articles/edit/page.sfc';
import renderer from 'sham-ui-test-helpers';
import { storage } from '../../../../../../src/storages/session';

it( 'renders correctly', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;
    DI.bind( 'router', {
        storage: {
            params: {
                id: 42
            }
        }
    } );

    DI.bind( 'store', {
        articleCategories: jest.fn().mockReturnValue(
            Promise.resolve( {
                categories: [
                    { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                    { 'id': '6', 'name': 'Название категории', 'slug': 'nazvanie-kategorii' }
                ],
                meta: {}
            } )
        ),
        articleTags: jest.fn().mockReturnValue(
            Promise.resolve( {
                tags: [
                    { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                    { 'id': '6', 'name': 'Название тега', 'slug': 'nazvanie-tega' }
                ],
                meta: {}
            } )
        ),
        articleDetail: jest.fn().mockReturnValue(
            Promise.resolve( {
                'title': 'Тест2',
                'slug': 'test2',
                'category_id': '6',
                'short_body': 'Короткое',
                'body': '<p>Текст</p><p><u>Текст</u></p>',
                'published_at': '2022-04-25T19:34:23.619+07:00',
                'tags': [ 'hello-world' ]
            } )
        ),
        api: {
            token: '123',
            baseURL: 'localhost'
        }
    } );
    const meta = renderer( RoutesArticlesEditPage, {}, {
        DI,
        directives: {
            ...directives,
            onFocusIn
        }
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'fail load article data', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;
    DI.bind( 'router', {
        storage: {
            params: {
                id: 42
            }
        }
    } );

    DI.bind( 'store', {
        articleCategories: jest.fn().mockReturnValue(
            Promise.resolve( {
                categories: [
                    { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                    { 'id': '6', 'name': 'Название категории', 'slug': 'nazvanie-kategorii' }
                ],
                meta: {}
            } )
        ),
        articleTags: jest.fn().mockReturnValue(
            Promise.resolve( {
                tags: [
                    { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                    { 'id': '6', 'name': 'Название тега', 'slug': 'nazvanie-tega' }
                ],
                meta: {}
            } )
        ),
        articleDetail: jest.fn().mockReturnValue(
            Promise.reject( {} )
        ),
        api: {
            token: '123',
            baseURL: 'localhost'
        }
    } );
    const meta = renderer( RoutesArticlesEditPage, {}, {
        DI,
        directives: {
            ...directives,
            onFocusIn
        }
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
