<template>
    <div class="container">
      <div class="section">
        <div class="row">
          <div class="col l6 m12" v-for="bucket in buckets" :key="bucket.name">
            <router-link :to="{ name: 'Bucket', params: { bucketName: bucket.name } }" style="color:black;">
              <div class="card">
                <div class="card-content">
                  <div class="row">
                    <div class="col">
                      <i class="material-icons large">folder_open</i>
                    </div>
                    <div class="col">
                      <span class="card-title">{{ bucket.name }}</span>
                      <p style="color:gray;">Created on {{ bucket.creationDate }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </router-link>
          </div>
          <p v-if="!buckets.length" style="text-align:center;margin-top:2em;color:gray;">No buckets yet</p>
        </div>
      </div>
      <div class="fixed-action-btn">
        <button type="button" class="btn-floating btn-large modal-trigger" data-target="modal-create-bucket" style="background-color: #1976D2;">
          <i class="material-icons large">add</i>
        </button>
      </div>
      <div id="modal-create-bucket" class="modal">
        <form id="create-bucket-form" @submit.prevent="createBucket">
          <div class="modal-content">
            <h4>Create Bucket</h4>
            <br>
            <div class="row">
              <div class="input-field col m6">
                <input id="name" type="text" name="name" placeholder="My Bucket" v-model="newBucketName">
                <label for="name">Name</label>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="modal-close waves-effect waves-green btn-flat">Cancel</button>
            <button type="submit" class="waves-effect waves-light btn" style="background-color: #1976D2;">Create</button>
          </div>
        </form>
      </div>
    </div>
  </template>
  
  <script>
  /* global M */
  import { API_ENDPOINTS } from '../config/api.js';
  
  // eslint-disable-next-line vue/multi-word-component-names
  export default {
    name: 'Buckets',
    data() {
      return {
        buckets: [],
        newBucketName: ''
      }
    },
    mounted() {
      this.fetchBuckets();
      M.Modal.init(document.querySelectorAll('.modal'));
    },
    methods: {
      async fetchBuckets() {
        try {
          const response = await fetch(API_ENDPOINTS.buckets);
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          const data = await response.json();
          // API 返回的是 { buckets: [...], count: 2 } 格式
          // 我们需要提取 buckets 数组
          this.buckets = data.buckets || data;
        } catch (error) {
          console.error('Error fetching buckets:', error);
          M.toast({ html: 'Failed to load buckets', classes: 'red' });
        }
      },
      async createBucket() {
        try {
          const response = await fetch(API_ENDPOINTS.buckets, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name: this.newBucketName, region: "us-east-1" })
          });
          
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          
          M.toast({ html: 'Bucket created successfully!', classes: 'green' });
          this.newBucketName = '';
          const modal = M.Modal.getInstance(document.getElementById('modal-create-bucket'));
          modal.close();
          // Refresh the buckets list
          this.fetchBuckets();
        } catch (error) {
          console.error('Error creating bucket:', error);
          M.toast({ html: 'Failed to create bucket', classes: 'red' });
        }
      }
    }
  }
  </script>