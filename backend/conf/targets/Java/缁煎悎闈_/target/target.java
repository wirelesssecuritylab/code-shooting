@Slf4j
@Service
public class RelationTableServiceImpl implements RelationTableService {
    @Resource
    private RelationTableMapper relationTableMapper;

    @Resource
    private EsService esService;

    private static final String STR_START_TIME = "startTime";
    private static final String STR_END_TIME = "endTime";
    private static final String STR_CLEAR_DATE = "clearDate";
    private static final String STR_STAT_DATE = "statDate";

    public Map<String, Date> buildTimeForProcessUse(Date curTime) {
        Date startTime = LocalUtils.getTodayZeroTime();
        Date endTime = curTime;
        Date statDate = curTime;
        long zeroThirtyMs = LocalUtils.calculateDate(startTime, Calendar.MINUTE, 30).getTime();
        long curMs = curTime.getTime();
        //如果定时任务时间在00:00:00-00:30:00执行
        if (curMs >= startTime.getTime() && curMs < zeroThirtyMs) {
            startTime = LocalUtils.getYesterdayZeroTime();
            endTime = LocalUtils.getTodayZeroTime();
            statDate = startTime;
        }
        Map<String, Date> dateMap = new HashMap<>();
        dateMap.put(STR_START_TIME, startTime);
        dateMap.put(STR_END_TIME, endTime);
        dateMap.put(STR_STAT_DATE, statDate);
        return dateMap;
    }

    @Override
    public ResponseEntity statUserProcessUseTime() throws IOException, ElasticsearchClientException {
        log.info("statUserProcessOnlineTime-->stat begin");
        Date curTime = new Date();
        Map<String, Date> dateMap = buildTimeForProcessUse(curTime);
        Date startTime = dateMap.get(STR_START_TIME);
        Date endTime = dateMap.get(STR_END_TIME);
        Date statDate = dateMap.get(STR_STAT_DATE);
        //时间格式:yyyy-MM-dd'T'HH:mm:ss.SSSXXX
        String startValue = LocalUtils.parseDateToStr(startTime, LocalConstant.DATE_FORMAT_TZ_MILLIS);
        String endValue = LocalUtils.parseDateToStr(endTime, LocalConstant.DATE_FORMAT_TZ_MILLIS);
        log.info("statUserProcessOnlineTime-->stat startValue={},endValue={}", startValue, endValue);
        //时间范围: startTime <= time < endTime
        //本版本暂时不按照节点分批查询
        String sumKey = LocalConstant.INDEX_KEY_USE_TIME;
        EsQueryDto.EsQueryDtoBuilder builder = new EsQueryDto.EsQueryDtoBuilder();
        EsQueryDto query = builder.rangeKey(LocalConstant.INDEX_KEY_CREATETIME).startValue(startValue).endValue(endValue).aggField(LocalConstant.INDEX_KEY_NODE_ID).secAggField(LocalConstant.INDEX_KEY_VMID).thirdAggField(LocalConstant.INDEX_KEY_PROCESS_NAME).includeLower(true).includeUpper(false).build();
        List<JSONObject> jsons = esService.statSumByThreeAgg(LocalConstant.INDEX_USER_PROCESS, query, new HashMap(), Arrays.asList(sumKey));
        int size = LocalUtils.getColSize(jsons);
        if (size == 0) {
            log.warn("statUserProcessOnlineTime-->es query result is empty");
            return ResponseEntity.failResponse("query nothing from es");
        }
        List<VmUserOrg> relations = relationTableMapper.queryVmUserOrg();
        if (LocalUtils.getColSize(relations) == 0) {
            log.error("statUserProcessOnlineTime-->vm user org data not exist");
            return ResponseEntity.failResponse("vm user org data not exist");
        }
        //分隔线
        String sep = LocalConstant.SEP_UNDERLINE;
        Map<String, VmUserOrg> vmUserOrgMap = new HashMap<>();
        relations.parallelStream().forEach(relation ->
        {
            log.info("will process {}", relation);
            vmUserOrgMap.put(relation.getNodeId() + sep + relation.getVmid(), relation);
        });
        List<SoftwareUseStat> processUseList = new ArrayList<>(size);
        String firstAggField = query.getAggField();
        String secAggField = query.getSecAggField();
        String thirdAggField = query.getThirdAggField();

        for (JSONObject json : jsons) {
            Integer nodeId = LocalUtil.parseObjToInt(json.get(firstAggField), 0);
            String vmid = (String) json.get(secAggField);
            String processName = (String) json.get(thirdAggField);
            int sumValue = LocalUtils.parseDoubleToInt((double) json.get(sumKey), 0);
            SoftwareUseStat stat = new SoftwareUseStat(nodeId, vmid, processName, sumValue);
            VmUserOrg vmUserOrg = vmUserOrgMap.get(nodeId + sep + vmid);
            if (null != vmUserOrg) {
                stat.setOrganizationId(vmUserOrg.getOrganizationId());
                stat.setOsType(vmUserOrg.getOsType());
                stat.setUserId(vmUserOrg.getUserId());
            }
            stat.setStatDate(statDate);
            processUseList.add(stat);
        }
        vmUserOrgMap.clear();
        //拆分list批量插入
        int total = batchInsertProcessUse(statDate, processUseList);
        log.info("statUserProcessOnlineTime-->stat end,size={}", size);

        return ResponseEntity.sucResponse("stat success", total);
    }

    @Transactional(rollbackFor = Exception.class, timeout = 600)
    public int batchInsertProcessUse(Date statDate, List<SoftwareUseStat> processUseList) {
        relationTableMapper.deleteProcessUseByDate(statDate);
        List<List<SoftwareUseStat>> splitedList = LocalUtils.splitList(processUseList, 2000);
        int total = 0;
        for (List<SoftwareUseStat> subList : splitedList) {
            int count = relationTableMapper.batchInsertProcessUse(subList);
            log.info("batchInsertProcessUse -->count={}", count);
            total += count;
        }
        return total;
    }

}