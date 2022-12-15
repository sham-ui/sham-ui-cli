import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesArticlesCommonFieldsEditor  from '../../../../../../../src/components/routes/articles/common/fields/editor.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'store', {
        api: {
            token: '123',
            baseURL: 'localhost'
        }
    } );
    const meta = renderer( RoutesArticlesCommonFieldsEditor, {}, {
        DI,
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
