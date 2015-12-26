/* jshint -W097,-W117 */
'use strict';

var React = require('react');
var _ = require('underscore');

class QueryEditBar extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleSaveChanges = this.handleSaveChanges.bind(this);
    this.handleSQLPreview = this.handleSQLPreview.bind(this);

    this.state = {
    };
  }

  handleSaveChanges() {
    console.log("handleSaveChanges ");
    // TODO: execute query
  }

  handleSQLPreview() {
    console.log("handleSQLPreview ");
    // TODO: show sqlpreview in modal
    // TODO: maybe move generateQuery from output to here?
    var query = this.props.generateQuery();
    console.log("Query Preview", query);
  }

  render() {
    return (
      <div id="query_edit_bar">
        <button className="discard_changes" onClick={this.props.onHandleDiscardChanges}>Discard Changes</button>
        <button className="sql_preview" onClick={this.handleSQLPreview.bind(this)}>SQL Preview</button>
        <button className="save_changes" onClick={this.handleSaveChanges.bind(this)}>Save Changes</button>
      </div>
    );
  }
}

module.exports = QueryEditBar;
