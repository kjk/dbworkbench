import React from 'react';
import ReactDOM from 'react-dom';

export class FiltererInput extends React.Component {
    onChange() {
        this.props.onFilter(ReactDOM.findDOMNode(this).value);
    }

    render() {
        return (
            <input type="text"
                className="reactable-filter-input"
                placeholder={this.props.placeholder}
                value={this.props.value}
                onKeyUp={this.onChange.bind(this)}
                onChange={this.onChange.bind(this)} />
        );
    }
};

export class Filterer extends React.Component {
    render() {
        return (
            <div className="reactable-filterer">
                <FiltererInput onFilter={this.props.onFilter}
                    value={this.props.value}
                    placeholder={this.props.placeholder}/>
            </div>
        );
    }
};

