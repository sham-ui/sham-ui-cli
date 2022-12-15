import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesArticlesTagsDetailForm  from '../../../../../../../src/components/routes/articles/tags/detail/form.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesArticlesTagsDetailForm, {}, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
