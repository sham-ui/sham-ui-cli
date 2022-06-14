import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import RoutesArticlesNewPage  from '../../../../../../src/components/routes/articles/new/page.sfc';
import renderer from 'sham-ui-test-helpers';
import { storage } from '../../../../../../src/storages/session';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;

    const getMock = jest.fn();
    DI.bind( 'store', {
        articleCategories: getMock.mockReturnValueOnce(
            Promise.resolve( {
                categories: [
                    { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                    { 'id': '6', 'name': 'Название категории', 'slug': 'nazvanie-kategorii' }
                ],
                meta: {}
            } )
        ),
        articleTags: getMock.mockReturnValueOnce(
            Promise.resolve( {
                tags: [
                    { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                    { 'id': '6', 'name': 'Название тега', 'slug': 'nazvanie-tega' }
                ],
                meta: {}
            } )
        )
    } );

    const meta = renderer( RoutesArticlesNewPage, {
        DI,
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
