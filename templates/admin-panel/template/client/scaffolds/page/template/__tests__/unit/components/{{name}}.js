import {{ classifiedName }}  from '{{testRelativePath}}';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( {{ classifiedName }}, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
