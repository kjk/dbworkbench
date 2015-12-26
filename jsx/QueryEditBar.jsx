/* jshint -W097,-W117 */
'use strict';

var React = require('react');
var _ = require('underscore');

var api = require('./api.js');

class QueryEditBar extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleSaveChanges = this.handleSaveChanges.bind(this);
    this.handleSQLPreview = this.handleSQLPreview.bind(this);

    // 1) Is there a way to move discard changes to here without using action?
    // 2) maybe move generateQuery from output to here?

    this.state = {
    };
  }

  handleSaveChanges() {
    console.log("handleSaveChanges ");
    // TODO: execute query
    // TODO: must support multiple queries for multiple rows changes
    // var query = this.props.generateQuery();
    // console.log("Executing query", query);

    // api.executeQuery(query);
  }

  handleSQLPreview() {
    console.log("handleSQLPreview ");
    // TODO: show sqlpreview in modal
    var query = this.props.generateQuery();
    console.log("Query Preview", query);
  }

  render() {
    return (
      <div id="query_edit_bar">
        <button className="discard_changes" onClick={this.props.onHandleDiscardChanges}>Discard Changes</button>
        <div className="row_number">{this.props.numberOfRowsEdited} edited rows</div>
        <button className="sql_preview" onClick={this.handleSQLPreview.bind(this)}>SQL Preview</button>
        <button className="save_changes" onClick={this.handleSaveChanges.bind(this)}>Save Changes</button>
      </div>
    );
  }
}

module.exports = QueryEditBar;
