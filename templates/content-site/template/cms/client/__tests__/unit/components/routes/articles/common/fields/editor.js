import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesArticlesCommonFieldsEditor  from '../../../../../../../src/components/routes/articles/common/fields/editor.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesArticlesCommonFieldsEditor, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
